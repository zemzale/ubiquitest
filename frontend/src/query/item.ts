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
        queryFn: () => JSON.parse(
            localStorage.getItem('todos') ?? '[]'
        ) as Item[]
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
    const ws = useWebsocket()
    const client = useQueryClient();
    return useMutation({
        mutationFn: completeItem(ws),
        onSuccess: () => {
            client.invalidateQueries({ queryKey: ['todos'] });
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

function completeItem(ws: WebSocket) {
    return async (itemId: string) => {
        ws.send(JSON.stringify({
            type: 'task_done',
            data: { id: itemId },
        }));

        const todos = JSON.parse(localStorage.getItem("todos") ?? "[]") as Item[];
        const updatedTodos = todos.map(todo =>
            todo.id === itemId ? { ...todo, completed: true } : todo
        );
        localStorage.setItem("todos", JSON.stringify(updatedTodos));
    }
}

