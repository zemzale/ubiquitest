import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { env } from '~/env'
import { Item } from './item'

export type User = {
    username: string;
    id: number;
}

export function useUser() {
    return useQuery({
        queryKey: ['user'],
        queryFn: () => {
            const userData = localStorage.getItem('user');
            if (!userData) return null;
            return JSON.parse(userData) as User;
        },
    });
}

export function useLogin() {
    const queryClient = useQueryClient();
    return useMutation({
        mutationFn: postLogin,
        onSuccess: (result: User) => {
            queryClient.invalidateQueries({ queryKey: ['user'] });
            localStorage.setItem('user', JSON.stringify(result));
            
            // Check if this is the first login (todos not yet fetched)
            const hasFetchedTodos = localStorage.getItem('hasFetchedTodos') === 'true';
            
            if (!hasFetchedTodos) {
                console.log('First login - fetching todos from server');
                
                // Pre-fetch tasks only on first login
                queryClient.prefetchQuery({
                    queryKey: ['todos'],
                    queryFn: async () => {
                        try {
                            const response = await fetch(`${env.NEXT_PUBLIC_API_URL}/todos`);
                            if (!response.ok) {
                                throw new Error('Failed to fetch tasks');
                            }
                            const todos = await response.json();
                            
                            // Store in localStorage as backup
                            localStorage.setItem('todos', JSON.stringify(todos));
                            
                            // Set the flag indicating we've fetched todos at least once
                            localStorage.setItem('hasFetchedTodos', 'true');
                            
                            console.log('Initial todos fetch complete and saved');
                            return todos;
                        } catch (error) {
                            console.error('Error pre-fetching tasks:', error);
                            return [];
                        }
                    }
                });
            } else {
                console.log('Subsequent login - not fetching todos from server');
            }
        },
    })
}

export function useUserById(userId: string | number | undefined) {
    return useQuery({
        queryKey: ['user', userId],
        queryFn: () => fetchUserById(userId),
        enabled: !!userId, // Only run the query if userId is provided
        staleTime: 5 * 60 * 1000, // Cache for 5 minutes
    });
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

async function fetchUserById(userId: string | number | undefined) {
    if (!userId) return null;
    
    return fetch(`${env.NEXT_PUBLIC_API_URL}/user/${userId}`)
        .then((res) => {
            if (!res.ok) {
                throw new Error('Failed to fetch user');
            }
            return res.json();
        });
}

