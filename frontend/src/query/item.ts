import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { v4 as uuidv4 } from 'uuid';
import { useWebsocket } from '~/ws/hook';

export type Item = {
    title: string;
    id: string,
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
    const ws = useWebsocket()
    const client = useQueryClient();
    return useMutation({
        mutationFn: postItem(ws),
        onSuccess: () => {
            client.invalidateQueries({ queryKey: ['todos'] });
        },
    })
}

type NewItem = Omit<Item, 'id'>;

function postItem(ws: WebSocket) {
    return async (body: NewItem) => {
        const id = uuidv4();
        const item: Item = {
            id: id,
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

