import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { v4 as uuidv4 } from 'uuid';
import { useWebsocket, EnhancedWebSocket } from '~/ws/hook';
import { User } from './user';
import { env } from '~/env';

export type Item = {
    title: string;
    id: string;
    completed?: boolean;
    created_by: number;
    parent_id?: string;
}

export function useItems() {
    return useQuery({
        queryKey: ['todos'],
        queryFn: async () => {
            try {
                // First try to fetch from server API
                const response = await fetch(`${env.NEXT_PUBLIC_API_URL}/todos`);
                if (!response.ok) {
                    throw new Error('Failed to fetch tasks from server');
                }
                const serverTodos = await response.json() as Item[];

                // Store in localStorage for offline access
                localStorage.setItem('todos', JSON.stringify(serverTodos));
                console.log('Fetched todos from server:', serverTodos);
                return serverTodos;
            } catch (error) {
                console.error('Error fetching from server, falling back to localStorage:', error);
                // Fallback to localStorage if server fetch fails
                const todosString = localStorage.getItem('todos');
                console.log('Fetching todos from localStorage:', todosString);
                return JSON.parse(todosString ?? '[]') as Item[];
            }
        },
        // Ensure the query refetches when invalidated
        staleTime: 0,
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

export function useUpdateTask() {
    const ws = useWebsocket();
    const client = useQueryClient();
    return useMutation({
        mutationFn: updateTask(ws, client),
        onSuccess: () => {
            client.invalidateQueries({ queryKey: ['todos'] });
            console.log('Invalidated todos query after task update');
        },
    });
}

export function useCompleteItem() {
    const ws = useWebsocket();
    const client = useQueryClient();
    return useMutation({
        mutationFn: completeItem(ws, client),
        onSuccess: () => {
            client.invalidateQueries({ queryKey: ['todos'] });
            console.log('Invalidated todos query after completion');
        },
    })
}

type NewItem = Omit<Item, 'id' | 'created_by'>;

function postItem(ws: EnhancedWebSocket, user: User) {
    return async (body: NewItem) => {
        const id = uuidv4();
        const item: Item = {
            id: id,
            completed: false,
            created_by: user.id,
            ...body,
        };

        // Remove undefined parent_id to avoid sending it in the payload
        if (item.parent_id === undefined) {
            delete item.parent_id;
        }

        // Use the enhanced send method which handles connection state
        ws.send(JSON.stringify({
            type: 'task_created',
            data: item,
        }));

        const todos = JSON.parse(localStorage.getItem("todos") ?? "[]") as Item[];
        todos.push(item);
        localStorage.setItem("todos", JSON.stringify(todos));
    }
}

function updateTask(ws: EnhancedWebSocket, queryClient: ReturnType<typeof useQueryClient>) {
    return async ({ id, changes }: { id: string, changes: Partial<Item> }) => {
        console.log('Updating task:', id, 'with changes:', changes);
        const todos = JSON.parse(localStorage.getItem("todos") ?? "[]") as Item[];
        const item = todos.find(todo => todo.id === id);

        if (!item) {
            console.error('Cannot update item, not found in localStorage:', id);
            return;
        }
        const updatedItem = { ...item, ...changes };

        // Use the enhanced send method which handles connection state
        ws.send(JSON.stringify({
            type: 'task_updated',
            data: updatedItem,
        }));
        console.log('Sent task_updated message with data:', updatedItem);

        const updatedTodos = todos.map(todo =>
            todo.id === id ? updatedItem : todo
        );
        localStorage.setItem("todos", JSON.stringify(updatedTodos));
        console.log('Updated todos in localStorage:', updatedTodos);

        // Explicitly invalidate the query to trigger a re-fetch
        queryClient.invalidateQueries({ queryKey: ['todos'] });
        console.log('Invalidated todos query from updateTask function');
    }
}

// Use the generic updateTask function for completing items
function completeItem(ws: EnhancedWebSocket, queryClient: ReturnType<typeof useQueryClient>) {
    const updateTaskFn = updateTask(ws, queryClient);

    return async (itemId: string) => {
        console.log('Completing task:', itemId);
        return updateTaskFn({
            id: itemId,
            changes: { completed: true }
        });
    }
}

/**
 * Type definition for an item with children (tree node)
 */
export type ItemWithChildren = Item & { children: ItemWithChildren[] };

/**
 * Organizes a flat list of items into a hierarchical tree structure.
 * Items with parent_id will be nested under their parent.
 */
export function organizeItemsIntoTree(items: Item[]): ItemWithChildren[] {
    const itemMap = new Map<string, ItemWithChildren>();
    const rootItems: ItemWithChildren[] = [];

    // First pass: Create a map of all items with empty children arrays
    items.forEach(item => {
        itemMap.set(item.id, { ...item, children: [] });
    });

    // Second pass: Organize items into tree structure
    items.forEach(item => {
        const enhancedItem = itemMap.get(item.id)!;

        if (item.parent_id && itemMap.has(item.parent_id)) {
            // This is a child item, add it to its parent's children
            const parent = itemMap.get(item.parent_id)!;
            parent.children.push(enhancedItem);
        } else {
            // This is a root item
            rootItems.push(enhancedItem);
        }
    });

    return rootItems;
}

