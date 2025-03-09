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
            const todos = JSON.parse(localStorage.getItem("todos") ?? "[]") as Item[];
            todos.push(JSON.parse(event.data));
            localStorage.setItem("todos", JSON.stringify(todos));

            queryClient.invalidateQueries({ queryKey: ['todos'] });
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
