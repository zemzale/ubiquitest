import { useQueryClient } from "@tanstack/react-query";
import { createContext, useContext, useMemo } from "react";
import { env } from "~/env";
import { Item } from "~/query/item";

const WebSocketContext = createContext<WebSocket | null>(null);

export function useCreateWebsocket(user: string) {
    const queryClient = useQueryClient();
    const ws = useMemo(() => {
        const ws = new WebSocket(`${env.NEXT_PUBLIC_API_URL}/ws/todos?user=${user}`);
        ws.addEventListener('open', () => {
            console.log('WebSocket opened');
        });

        ws.addEventListener("message", (event) => {
            console.log('WebSocket message received:', event.data);

            try {
                const message = JSON.parse(event.data);
                const todos = JSON.parse(localStorage.getItem("todos") ?? "[]") as Item[];

                if (message.type === 'task_created') {
                    const exists = todos.some(todo => todo.id === message.data.id);
                    if (exists) {
                        console.warn('Received task_created for existing task:', message.data.id);
                    }
                    todos.push(message.data);
                    localStorage.setItem("todos", JSON.stringify(todos));
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

        return ws;
    }, [user]);

    return ws;
}

export function useWebsocket() {
    const ws = useContext(WebSocketContext);
    if (!ws) {
        throw new Error('useWebsocket must be used within a WebSocketProvider');
    }

    return ws;
}

export function WebSocketProvider({ value, children }: { value: WebSocket, children: React.ReactNode }) {

    return (
        <WebSocketContext.Provider value={value} >
            {children}
        </WebSocketContext.Provider>
    );
}
