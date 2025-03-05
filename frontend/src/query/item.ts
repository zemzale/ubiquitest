import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { env } from '~/env'
import { v4 as uuidv4 } from 'uuid';
import { useMemo } from 'react';
import { useUser } from './user';

type Item = {
    title: string;
    id: string,
}


export function useWebsocket() {
    const ws = useMemo(() => {
        const ws = new WebSocket(`${env.NEXT_PUBLIC_API_URL}/ws/todos`);
        ws.addEventListener('open', () => {
            console.log('open');
        });

        ws.addEventListener('message', (event) => {
            console.log("Message from server:", event.data);
        });
        return ws;
    }, []);
    return ws;
}

export function useItems() {
    return useQuery({
        queryKey: ['todos'],
        queryFn: () => fetch(`${env.NEXT_PUBLIC_API_URL}/todos`)
            .then((res) => res.json())
            .then((data) => data as Item[]),
    });
}

export function useAddItem() {
    const queryClient = useQueryClient();
    const ws = useWebsocket();
    const user = useUser();
    return useMutation({
        mutationFn: postItem(ws, user),
        onMutate: () => {
            queryClient.invalidateQueries({ queryKey: ['todos'] });
        },
    })
}

type NewItem = Omit<Item, 'id'>;

function postItem(ws: WebSocket, user: any) {
    return async (body: NewItem) => {
        const id = uuidv4();
        const item: Item = {
            id: id,
            ...body,
        };

        if (user.data) {
            ws.send(JSON.stringify({
                type: 'task_created',
                data: {
                    id: id,
                    title: body.title,
                    created_by: user.data.id,
                },
            }));
        }


        return fetch(`${env.NEXT_PUBLIC_API_URL}/todos`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(item),
        })

            .then((res) => res.json())
    }
}

