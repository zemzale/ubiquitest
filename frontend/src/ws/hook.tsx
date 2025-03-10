import { useQueryClient } from "@tanstack/react-query";
import { createContext, useContext, useEffect, useRef, useState } from "react";
import { env } from "~/env";
import { Item } from "~/query/item";

// WebSocket status type - making the connection state more explicit
export type WebSocketStatus = 'connecting' | 'connected' | 'disconnected' | 'reconnecting';

// Enhanced WebSocket with reconnection capabilities and status
export interface EnhancedWebSocket {
    socket: WebSocket | null;
    status: WebSocketStatus;
    reconnect: () => void;
    send: (data: string) => void;
}

// Enhanced context with more information about the WebSocket connection
const WebSocketContext = createContext<EnhancedWebSocket | null>(null);

// Default retry configuration
const DEFAULT_RECONNECT_DELAY_MS = 10000;
const MAX_RECONNECT_DELAY_MS = 30000;
const RECONNECT_BACKOFF_FACTOR = 3;
const MAX_RECONNECT_ATTEMPTS = 3;

// Ping configuration (send a ping every second to keep connection alive)
const PING_INTERVAL_MS = 1000;

export function useCreateWebsocket(user: string): EnhancedWebSocket {
    const queryClient = useQueryClient();
    const [status, setStatus] = useState<WebSocketStatus>('connecting');
    const [socket, setSocket] = useState<WebSocket | null>(null);

    // Use refs to maintain state across re-renders and in event handlers
    const reconnectTimeoutRef = useRef<NodeJS.Timeout | null>(null);
    const pingIntervalRef = useRef<NodeJS.Timeout | null>(null);
    const reconnectAttemptsRef = useRef(0);
    const currentDelayRef = useRef(DEFAULT_RECONNECT_DELAY_MS);

    // Function to create a new socket connection
    const createSocket = () => {
        try {
            // Clear any existing reconnect timeout
            if (reconnectTimeoutRef.current) {
                clearTimeout(reconnectTimeoutRef.current);
                reconnectTimeoutRef.current = null;
            }

            // Update status to connecting
            setStatus('connecting');

            // Create a new WebSocket connection
            const newWs = new WebSocket(`${env.NEXT_PUBLIC_API_URL}/ws/todos?user=${user}`);

            // Configure event handlers
            newWs.addEventListener('open', () => {
                console.log('WebSocket connected successfully');
                setStatus('connected');
                reconnectAttemptsRef.current = 0;
                currentDelayRef.current = DEFAULT_RECONNECT_DELAY_MS;
                
                // Setup ping interval to keep connection alive
                if (pingIntervalRef.current) {
                    clearInterval(pingIntervalRef.current);
                }
                
                // Keep track of ping count to reduce logging
                let pingCount = 0;
                
                pingIntervalRef.current = setInterval(() => {
                    if (newWs.readyState === WebSocket.OPEN) {
                        // Send a ping message to keep the connection alive
                        newWs.send(JSON.stringify({ type: 'ping', timestamp: Date.now() }));
                        
                    }
                }, PING_INTERVAL_MS);
            });

            newWs.addEventListener('close', (event) => {
                console.log(`WebSocket closed with code ${event.code}, reason: ${event.reason}`);
                setStatus('disconnected');
                
                // Clear the ping interval
                if (pingIntervalRef.current) {
                    clearInterval(pingIntervalRef.current);
                    pingIntervalRef.current = null;
                }
                
                // Only attempt reconnects for unexpected disconnects
                if (event.code !== 1000) { // 1000 is normal closure
                    scheduleReconnect();
                }
                setSocket(null);
            });

            newWs.addEventListener('error', (error) => {
                console.error('WebSocket error:', error);
                // Error event is followed by close event, no need to set status or reconnect here
            });

            newWs.addEventListener("message", (event) => {
                console.log('WebSocket message received:', event.data);

                try {
                    const message = JSON.parse(event.data);
                    const todos = JSON.parse(localStorage.getItem("todos") ?? "[]") as Item[];

                    if (message.type === 'task_created') {
                        const exists = todos.some(todo => todo.id === message.data.id);
                        if (exists) {
                            console.warn('Received task_created for existing task:', message.data.id);
                        }

                        // Handle both with and without parent_id
                        const newTask = message.data;
                        todos.push(newTask);
                        localStorage.setItem("todos", JSON.stringify(todos));
                        console.log('Task created with ID:', newTask.id,
                            newTask.parent_id ? `as a subtask of ${newTask.parent_id}` : 'as a top-level task');
                    } else if (message.type === 'task_updated') {
                        const taskId = message.data.id;

                        if (!taskId) {
                            console.error(`Missing task ID in ${message.type} message:`, message);
                            return;
                        }

                        const taskExists = todos.some(todo => todo.id === taskId);

                        if (!taskExists) {
                            console.warn(`Received ${message.type} for unknown task:`, taskId);
                        }

                        const updatedTodos = todos.map(todo => {
                            if (todo.id === taskId) {
                                return { ...todo, ...message.data };
                            }
                            return todo;
                        });

                        localStorage.setItem("todos", JSON.stringify(updatedTodos));
                        console.log('Updated todos saved to localStorage');
                    } else {
                        console.log('Unknown message type:', message.type);
                    }

                    // Invalidate the todos query to trigger a re-fetch
                    queryClient.invalidateQueries({ queryKey: ['todos'] });
                } catch (error) {
                    console.error('Error processing WebSocket message:', error);
                }
            });

            // Update state with the new socket
            setSocket(newWs);
            return newWs;
        } catch (err) {
            console.error('Failed to create WebSocket connection:', err);
            setStatus('disconnected');
            scheduleReconnect();
            return null;
        }
    };

    // Schedule a reconnection attempt with exponential backoff
    const scheduleReconnect = () => {
        if (reconnectAttemptsRef.current >= MAX_RECONNECT_ATTEMPTS) {
            console.log(`Maximum reconnection attempts (${MAX_RECONNECT_ATTEMPTS}) reached, giving up`);
            return;
        }

        reconnectAttemptsRef.current += 1;

        // Calculate delay with exponential backoff
        const delay = Math.min(currentDelayRef.current, MAX_RECONNECT_DELAY_MS);

        console.log(`Scheduling reconnect attempt ${reconnectAttemptsRef.current} in ${delay}ms`);
        setStatus('reconnecting');

        // Schedule the reconnect attempt
        reconnectTimeoutRef.current = setTimeout(() => {
            console.log(`Attempting to reconnect (attempt ${reconnectAttemptsRef.current})`);
            createSocket();

            // Increase delay for next attempt
            currentDelayRef.current = Math.min(
                currentDelayRef.current * RECONNECT_BACKOFF_FACTOR,
                MAX_RECONNECT_DELAY_MS
            );
        }, delay);
    };

    // Function to manually trigger a reconnection
    const reconnect = () => {
        console.log('Manual reconnection initiated');

        // Close existing socket if it's open
        if (socket && (socket.readyState === WebSocket.OPEN || socket.readyState === WebSocket.CONNECTING)) {
            socket.close();
        }

        // Reset reconnection attempts for manual reconnections
        reconnectAttemptsRef.current = 0;
        currentDelayRef.current = DEFAULT_RECONNECT_DELAY_MS;

        // Create a new socket immediately
        createSocket();
    };

    // Initialize WebSocket connection when the component mounts or user changes
    useEffect(() => {
        const ws = createSocket();

        // Cleanup function to close the socket and clear timeouts when unmounting
        return () => {
            if (socket && socket.readyState === WebSocket.OPEN) {
                socket.close(1000, 'Component unmounted');
            }

            if (reconnectTimeoutRef.current) {
                clearTimeout(reconnectTimeoutRef.current);
            }
            if (pingIntervalRef.current) {
                clearInterval(pingIntervalRef.current);
                pingIntervalRef.current = null;
            }
        };
    }, [user]); // Recreate socket when user changes

    // Create a wrapper for the send method that handles socket state
    const send = (data: string) => {
        if (socket && socket.readyState === WebSocket.OPEN) {
            socket.send(data);
        } else {
            console.warn('Attempted to send message when socket is not open:', data);
            // Store the message to be sent when reconnected
            // (Optional enhancement: could implement a message queue here)
        }
    };

    // Return the enhanced WebSocket object
    return {
        socket,
        status,
        reconnect,
        send
    };
}

export function useWebsocket() {
    const ws = useContext(WebSocketContext);
    if (!ws) {
        throw new Error('useWebsocket must be used within a WebSocketProvider');
    }
    return ws;
}

export function WebSocketProvider({ value, children }: { value: EnhancedWebSocket, children: React.ReactNode }) {
    return (
        <WebSocketContext.Provider value={value}>
            {children}
        </WebSocketContext.Provider>
    );
}
