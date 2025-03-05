import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { env } from '~/env'
import { v4 as uuidv4 } from 'uuid';

type Item = {
    title: string;
    id: string,
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
    return useMutation({
        mutationFn: postItem,
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ['todos'] });
        },
    })
}

type NewItem = Omit<Item, 'id'>;

async function postItem(body: NewItem) {
    const id = uuidv4();
    const item: Item = {
        id: id,
        ...body,
    };

    return fetch(`${env.NEXT_PUBLIC_API_URL}/todos`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify(item),
    })
        .then((res) => res.json())
}

