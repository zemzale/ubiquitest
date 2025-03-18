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
        queryKey: ['tasks'],
        queryFn: async () => {
            const hasFetchedTasks = localStorage.getItem('hasFetchedTasks') === 'true';
            const tasksString = localStorage.getItem('tasks');

            // If we've already fetched tasks and have them in localStorage, use that
            if (hasFetchedTasks && tasksString) {
                console.log('Using cached tasks from localStorage');
                return JSON.parse(tasksString) as Item[];
            }

            // Otherwise fetch from server (first login or explicit refresh)
            try {
                console.log('Fetching tasks from server');
                const response = await fetch(`${env.NEXT_PUBLIC_API_URL}/tasks`);
                if (!response.ok) {
                    throw new Error('Failed to fetch tasks from server');
                }
                const serverTasks = await response.json() as Item[];

                // Store in localStorage for offline access
                localStorage.setItem('tasks', JSON.stringify(serverTasks));

                // Set the flag indicating we've fetched tasks
                localStorage.setItem('hasFetchedTasks', 'true');

                console.log('Fetched tasks from server:', serverTasks);
                return serverTasks;
            } catch (error) {
                console.error('Error fetching from server, falling back to localStorage:', error);

                // If we have data in localStorage, use it as fallback
                if (tasksString) {
                    console.log('Falling back to localStorage tasks');
                    return JSON.parse(tasksString) as Item[];
                }

                // Otherwise return empty array
                console.log('No tasks available, returning empty array');
                return [] as Item[];
            }
        },
        // Keep cached data longer, since we're relying on localStorage
        staleTime: 5 * 60 * 1000, // 5 minutes
    });
}

export function useAddItem() {
    const ws = useWebsocket();
    const client = useQueryClient();
    const user = JSON.parse(localStorage.getItem('user') || '{}');

    return useMutation({
        mutationFn: postItem(ws, user),
        onSuccess: () => {
            client.invalidateQueries({ queryKey: ['tasks'] });
        },
    })
}

export function useUpdateTask() {
    const ws = useWebsocket();
    const client = useQueryClient();
    return useMutation({
        mutationFn: updateTask(ws, client),
        onSuccess: () => {
            client.invalidateQueries({ queryKey: ['tasks'] });
            console.log('Invalidated tasks query after task update');
        },
    });
}

export function useCompleteItem() {
    const ws = useWebsocket();
    const client = useQueryClient();
    return useMutation({
        mutationFn: completeItem(ws, client),
        onSuccess: () => {
            client.invalidateQueries({ queryKey: ['tasks'] });
            console.log('Invalidated tasks query after completion');
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

        const tasks = JSON.parse(localStorage.getItem("tasks") ?? "[]") as Item[];
        tasks.push(item);
        localStorage.setItem("tasks", JSON.stringify(tasks));
    }
}

function updateTask(ws: EnhancedWebSocket, queryClient: ReturnType<typeof useQueryClient>) {
    return async ({ id, changes }: { id: string, changes: Partial<Item> }) => {
        console.log('Updating task:', id, 'with changes:', changes);
        const tasks = JSON.parse(localStorage.getItem("tasks") ?? "[]") as Item[];
        const item = tasks.find(todo => todo.id === id);

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

        const updatedTasks = tasks.map(todo =>
            todo.id === id ? updatedItem : todo
        );
        localStorage.setItem("tasks", JSON.stringify(updatedTasks));
        console.log('Updated tasks in localStorage:', updatedTasks);

        // Explicitly invalidate the query to trigger a re-fetch
        queryClient.invalidateQueries({ queryKey: ['tasks'] });
        console.log('Invalidated tasks query from updateTask function');
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
