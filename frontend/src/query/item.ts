import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { v4 as uuidv4 } from 'uuid';
import { useWebsocket } from '~/ws/hook';
import { User } from './user';

export type Item = {
    title: string;
    id: string;
    completed?: boolean;
    created_by: number;
}


export function useItems() {
    return useQuery({
        queryKey: ['todos'],
        queryFn: () => {
            const todosString = localStorage.getItem('todos');
            console.log('Fetching todos from localStorage:', todosString);
            return JSON.parse(todosString ?? '[]') as Item[];
        },
        // Ensure the query refetches when invalidated
        staleTime: 0,
        // Don't cache the result for long
        cacheTime: 5000,
    });
}

export function useAddItem() {
    const ws = useWebsocket();
    const client = useQueryClient();
    const user = JSON.parse(localStorage.getItem('user') || '{}');

    return useMutation({
        mutationFn: postItem(ws, user),
        onSuccess: () => {
            client.invalidateQueries({ queryKey: ['todos'] });
        },
    })
}

export function useCompleteItem() {
    const ws = useWebsocket();
    const client = useQueryClient();
    return useMutation({
        mutationFn: completeItem(ws, client),
        onSuccess: () => {
            // Explicitly invalidate the todos query to force a refresh
            client.invalidateQueries({ queryKey: ['todos'] });
            console.log('Invalidated todos query after completion');
        },
    })
}

type NewItem = Omit<Item, 'id' | 'created_by'>;

function postItem(ws: WebSocket, user: User) {
    return async (body: NewItem) => {
        const id = uuidv4();
        const item: Item = {
            id: id,
            completed: false,
            created_by: user.id,
            ...body,
        };

        ws.send(JSON.stringify({
            type: 'task_created',
            data: item,
        }));

        const todos = JSON.parse(localStorage.getItem("todos") ?? "[]") as Item[];
        todos.push(item);
        localStorage.setItem("todos", JSON.stringify(todos));
    }
}

function completeItem(ws: WebSocket, queryClient: ReturnType<typeof useQueryClient>) {
    return async (itemId: string) => {
        console.log('Sending task_done WebSocket message for item:', itemId);
        
        // Get the item data from localStorage to include in the message
        const todos = JSON.parse(localStorage.getItem("todos") ?? "[]") as Item[];
        const item = todos.find(todo => todo.id === itemId);
        
        if (!item) {
            console.error('Cannot complete item, not found in localStorage:', itemId);
            return;
        }
        
        // Update the item to mark it as completed
        const completedItem = { ...item, completed: true };
        
        // Send a complete message with the full item data
        ws.send(JSON.stringify({
            type: 'task_done',
            data: completedItem,
        }));
        console.log('Sent task_done message with data:', completedItem);

        // Update the item in localStorage
        const updatedTodos = todos.map(todo =>
            todo.id === itemId ? completedItem : todo
        );
        localStorage.setItem("todos", JSON.stringify(updatedTodos));
        console.log('Updated todos in localStorage:', updatedTodos);
        
        // Explicitly invalidate the query to trigger a re-fetch
        queryClient.invalidateQueries({ queryKey: ['todos'] });
        console.log('Invalidated todos query from completeItem function');
    }
}

