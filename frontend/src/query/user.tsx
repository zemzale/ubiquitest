import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { env } from '~/env'

export type User = {
    username: string;
    id: number;
}

export function useUser() {
    return useQuery({
        queryKey: ['user'],
        queryFn: () => JSON.parse(localStorage.getItem('user') as string) as User,
    });
}

export function useLogin() {
    const queryClient = useQueryClient();
    return useMutation({
        mutationFn: postLogin,
        onSuccess: (result: User) => {
            queryClient.invalidateQueries({ queryKey: ['user'] });
            localStorage.setItem('user', JSON.stringify(result));
        },
    })
}

async function postLogin(body: { username: string }) {
    return fetch(`${env.NEXT_PUBLIC_API_URL}/login`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify(body),
    })
        .then((res) => res.json())
}

