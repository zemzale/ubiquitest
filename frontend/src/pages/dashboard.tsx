import Head from "next/head";
import Link from "next/link";
import { useRouter } from "next/router";
import { useEffect, useRef, useState } from "react";
import { useQueryClient } from "@tanstack/react-query";
import { useAddItem, useCompleteItem, useItems, ItemWithChildren, organizeItemsIntoTree } from "~/query/item";
import { User, useUser, useUserById } from "~/query/user";
import { useCreateWebsocket, useWebsocket, WebSocketProvider } from "~/ws/hook";

export default function Dashboard() {
    const query = useUser();

    if (query.isLoading) {
        return <Loading />
    }

    if (query.isError) {
        return <Error />;
    }

    if (!query.data) {
        return <Error />;
    }


    return (
        <>
            <Head>
                <title>Dashboard | Ubiquitodo</title>
            </Head>
            <Page user={query.data} />
        </>
    );
}

function Page({ user }: { user: User }) {
    const ws = useCreateWebsocket(user.username);
    const [showModal, setShowModal] = useState(false);

    return (
        <>
            <WebSocketProvider value={ws}>
                <Navbar user={user} onAddTask={() => setShowModal(true)} />
                <main className="flex min-h-screen flex-col items-center pt-16 bg-gray-50">
                    <ListItems />
                </main>
                {showModal && <AddItemModal onClose={() => setShowModal(false)} />}
            </WebSocketProvider>
        </>
    );
}


function Loading() {
    return (
        <>
            <nav className="fixed top-0 left-0 right-0 bg-white shadow-md p-4 z-10">
                <div className="container mx-auto flex justify-between items-center">
                    <div className="font-bold text-xl">Ubiquitodo</div>
                </div>
            </nav>
            <main className="flex min-h-screen flex-col items-center justify-center pt-16">
                <p className="text-2xl mb-8">
                    Loading...
                </p>
            </main>
        </>
    );
}

function ListItems() {
    const query = useItems();
    const completeMutation = useCompleteItem();
    const queryClient = useQueryClient();

    if (!query.data) {
        return <></>;
    }

    const handleComplete = (id: string) => {
        completeMutation.mutate(id);
    };

    // Function to force refresh tasks from server
    const handleRefresh = () => {
        // Remove the flag to trigger a server fetch
        localStorage.removeItem('hasFetchedTasks');
        // Invalidate the query to trigger a refetch
        queryClient.invalidateQueries({ queryKey: ['tasks'] });
    };

    // Sort items: incomplete tasks first, then completed tasks
    const sortedItems = [...query.data].sort((a, b) => {
        if (a.completed && !b.completed) return 1;
        if (!a.completed && b.completed) return -1;
        return 0;
    });

    // Group items into active and completed
    const activeItems = sortedItems.filter(item => !item.completed);
    const completedItems = sortedItems.filter(item => item.completed);
    const hasCompletedItems = completedItems.length > 0;

    // Organize active items into a tree structure
    const activeItemsTree = organizeItemsIntoTree(activeItems);
    const completedItemsTree = organizeItemsIntoTree(completedItems);

    return <>
        <div className="w-full max-w-6xl mb-6 mt-6 px-4">
            <div className="flex justify-between items-center mb-6">
                <h2 className="text-2xl font-semibold">Task Dashboard</h2>
                <button
                    onClick={handleRefresh}
                    disabled={query.isLoading}
                    className="text-blue-500 hover:text-blue-700 text-sm flex items-center"
                >
                    <svg xmlns="http://www.w3.org/2000/svg" className="h-4 w-4 mr-1" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
                    </svg>
                    {query.isLoading ? 'Refreshing...' : 'Refresh'}
                </button>
            </div>

            <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                {/* Active Tasks Column */}
                <div className="bg-white rounded-lg shadow-md p-4">
                    <h3 className="text-lg font-semibold mb-4 text-blue-600 border-b pb-2">Upcoming Tasks</h3>
                    {activeItemsTree.length === 0 ? (
                        <p className="text-gray-500 italic text-center py-8">No upcoming tasks</p>
                    ) : (
                        <ul className="space-y-3 max-h-[calc(100vh-240px)] overflow-y-auto pr-2">
                            {activeItemsTree.map((item) => (
                                <TaskItem
                                    key={item.id}
                                    item={item}
                                    level={0}
                                    onComplete={handleComplete}
                                    isCompleted={false}
                                />
                            ))}
                        </ul>
                    )}
                </div>

                {/* Completed Tasks Column */}
                <div className="bg-white rounded-lg shadow-md p-4">
                    <h3 className="text-lg font-semibold mb-4 text-green-600 border-b pb-2">Completed Tasks</h3>
                    {!hasCompletedItems ? (
                        <p className="text-gray-500 italic text-center py-8">No completed tasks</p>
                    ) : (
                        <ul className="space-y-3 max-h-[calc(100vh-240px)] overflow-y-auto pr-2">
                            {completedItemsTree.map((item) => (
                                <TaskItem
                                    key={item.id}
                                    item={item}
                                    level={0}
                                    onComplete={handleComplete}
                                    isCompleted={true}
                                />
                            ))}
                        </ul>
                    )}
                </div>
            </div>
        </div>
    </>;
}

// Enhanced component that renders a task and its subtasks recursively
function TaskItem({
    item,
    level,
    onComplete,
    isCompleted
}: {
    item: ItemWithChildren,
    level: number,
    onComplete: (id: string) => void,
    isCompleted: boolean
}) {
    const hasChildren = item.children && item.children.length > 0;
    const [showAddSubtask, setShowAddSubtask] = useState(false);

    // Calculate styling based on nesting level
    const isSubtask = level > 0;

    // For subtasks, make them visually distinct with smaller size and darker background
    const subtaskClasses = isSubtask ?
        'mt-1 border-l-2 border-l-blue-300 pl-3 transform scale-95 origin-top-left' :
        '';

    return (
        <>
            <li key={item.id} className={`${subtaskClasses} ${isSubtask ? 'my-2' : 'mb-3'}`}>
                <div className={`${isCompleted ? 'bg-gray-50' : 'bg-white'} 
                    shadow-sm rounded-lg ${isSubtask ? 'p-2.5' : 'p-4'} hover:shadow-md transition-shadow duration-200 
                    border ${isCompleted ? 'border-gray-200' : 'border-gray-100'}`}>
                    <div className="flex items-center justify-between">
                        <div className="flex items-center flex-grow">
                            {isCompleted ? (
                                <div className={`${isSubtask ? 'w-4 h-4' : 'w-5 h-5'} rounded-full bg-green-100 border border-green-300 flex items-center justify-center mr-3 flex-shrink-0`}>
                                    <svg className={`${isSubtask ? 'w-2.5 h-2.5' : 'w-3 h-3'} text-green-600`} fill="currentColor" viewBox="0 0 20 20">
                                        <path fillRule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clipRule="evenodd" />
                                    </svg>
                                </div>
                            ) : (
                                <div className={`${isSubtask ? 'w-4 h-4' : 'w-5 h-5'} rounded-full bg-blue-50 border border-blue-200 mr-3 flex-shrink-0`}></div>
                            )}
                            <div className="min-w-0 flex-grow">
                                {isSubtask && (
                                    <div className="text-xs font-semibold uppercase mb-0.5">Subtask</div>
                                )}
                                <h3 className={`${isSubtask ? 'text-sm' : 'text-base'} font-medium 
                                    ${isCompleted ? 'text-gray-500 line-through' : 'text-gray-800'} 
                                    break-words`}>
                                    {item.title}
                                </h3>
                                {item.created_by && (
                                    <TaskCreator userId={item.created_by} isSubtask={isSubtask} />
                                )}
                                {hasChildren && (
                                    <div className="mt-1 text-xs text-blue-500">
                                        {item.children.length} {item.children.length === 1 ? 'subtask' : 'subtasks'}
                                    </div>
                                )}
                            </div>
                        </div>
                        <div className="flex space-x-1 ml-2 flex-shrink-0">
                            {!isCompleted && (
                                <button
                                    onClick={() => setShowAddSubtask(true)}
                                    className={`p-1 text-blue-500 hover:bg-blue-50 rounded-full ${isSubtask ? 'scale-90' : ''}`}
                                    title="Add subtask"
                                >
                                    <svg xmlns="http://www.w3.org/2000/svg" className={`${isSubtask ? 'h-4 w-4' : 'h-5 w-5'}`} viewBox="0 0 20 20" fill="currentColor">
                                        <path fillRule="evenodd" d="M10 5a1 1 0 011 1v3h3a1 1 0 110 2h-3v3a1 1 0 11-2 0v-3H6a1 1 0 110-2h3V6a1 1 0 011-1z" clipRule="evenodd" />
                                    </svg>
                                </button>
                            )}
                            {!isCompleted && (
                                <button
                                    onClick={() => onComplete(item.id)}
                                    className={`p-1 text-green-500 hover:bg-green-50 rounded-full ${isSubtask ? 'scale-90' : ''}`}
                                    title="Mark as complete"
                                >
                                    <svg xmlns="http://www.w3.org/2000/svg" className={`${isSubtask ? 'h-4 w-4' : 'h-5 w-5'}`} viewBox="0 0 20 20" fill="currentColor">
                                        <path fillRule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clipRule="evenodd" />
                                    </svg>
                                </button>
                            )}
                        </div>
                    </div>
                </div>
            </li>

            {/* Modal for adding subtask */}
            {showAddSubtask && (
                <AddItemModal
                    onClose={() => setShowAddSubtask(false)}
                    parentId={item.id}
                />
            )}

            {/* Render children if any */}
            {hasChildren && (
                <ul className="pl-6 mt-1 space-y-1 border-l-2 border-blue-200">
                    {item.children.map((child: any) => (
                        <TaskItem
                            key={child.id}
                            item={child}
                            level={level + 1}
                            onComplete={onComplete}
                            isCompleted={isCompleted}
                        />
                    ))}
                </ul>
            )}
        </>
    );
}

function TaskCreator({ userId, isSubtask = false }: { userId: number, isSubtask?: boolean }) {
    const { data: user, isLoading, isError } = useUserById(userId);

    // Style based on whether this is a subtask
    const textColorClass = isSubtask ? 'text-blue-400' : 'text-gray-500';

    if (isLoading) {
        return <p className={`text-xs ${textColorClass}`}>Loading creator...</p>;
    }

    if (isError || !user) {
        return <p className={`text-xs ${textColorClass}`}>Unknown creator</p>;
    }

    return <p className={`text-xs ${textColorClass}`}>Created by: {user.username}</p>;
}

function Error() {
    return <>
        <h1 className="text-6xl font-bold mb-8">Error</h1>
        <p className="text-2xl mb-8">
            Something went wrong
        </p>
    </>;
}

function Navbar({ user, onAddTask }: { user: User, onAddTask: () => void }) {
    const router = useRouter();
    const { status, reconnect } = useWebsocket();

    const handleLogout = () => {
        // Remove user data and tasks
        localStorage.removeItem('user');
        localStorage.removeItem('tasks');

        // Also remove the flag that tracks todo fetching
        // This means the user will get a fresh todo fetch on next login
        localStorage.removeItem('hasFetchedTasks');

        router.push('/');
    };

    // Get status color based on connection state
    const getStatusColor = () => {
        switch (status) {
            case 'connected': return 'bg-green-500';
            case 'connecting': return 'bg-yellow-500';
            case 'reconnecting': return 'bg-yellow-500 animate-pulse';
            case 'disconnected': return 'bg-red-500';
            default: return 'bg-gray-500';
        }
    };

    // Get status label based on connection state
    const getStatusLabel = () => {
        switch (status) {
            case 'connected': return 'Connected';
            case 'connecting': return 'Connecting...';
            case 'reconnecting': return 'Reconnecting...';
            case 'disconnected': return 'Disconnected';
            default: return 'Unknown';
        }
    };

    return (
        <nav className="fixed top-0 left-0 right-0 bg-white shadow-md p-4 z-10">
            <div className="container mx-auto max-w-6xl flex justify-between items-center px-4">
                <div className="flex items-center">
                    <div className="font-bold text-xl text-blue-600">Ubiquitodo</div>
                    <div className="flex items-center ml-6">
                        <div className={`h-3 w-3 rounded-full ${getStatusColor()} mr-2`}></div>
                        <span className="text-xs text-gray-600">{getStatusLabel()}</span>
                        {status !== 'connected' && (
                            <button
                                onClick={reconnect}
                                className="ml-2 text-xs text-blue-600 hover:text-blue-800 underline"
                                title="Reconnect WebSocket"
                            >
                                Reconnect
                            </button>
                        )}
                    </div>
                </div>
                <div className="flex items-center">
                    <button
                        onClick={onAddTask}
                        className="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded text-sm mr-4 flex items-center"
                    >
                        <svg xmlns="http://www.w3.org/2000/svg" className="h-4 w-4 mr-1" viewBox="0 0 20 20" fill="currentColor">
                            <path fillRule="evenodd" d="M10 5a1 1 0 011 1v3h3a1 1 0 110 2h-3v3a1 1 0 11-2 0v-3H6a1 1 0 110-2h3V6a1 1 0 011-1z" clipRule="evenodd" />
                        </svg>
                        New Task
                    </button>
                    <div className="flex items-center bg-gray-100 rounded-full px-3 py-1 mr-3">
                        <span className="text-sm mr-2">
                            <span className="text-gray-600">@</span><span className="font-semibold text-gray-800">{user.username}</span>
                        </span>
                        <button
                            onClick={handleLogout}
                            className="text-gray-500 hover:text-red-500"
                            title="Logout"
                        >
                            <svg xmlns="http://www.w3.org/2000/svg" className="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
                                <path fillRule="evenodd" d="M3 3a1 1 0 00-1 1v12a1 1 0 001 1h12a1 1 0 001-1V7.414l-5-5H3zm9 5a1 1 0 10-2 0v4a1 1 0 102 0V8zm-2-7a1 1 0 00-1 1v1a1 1 0 102 0V2a1 1 0 00-1-1z" clipRule="evenodd" />
                            </svg>
                        </button>
                    </div>
                </div>
            </div>
        </nav>
    );
}


function AddItemModal({ onClose, parentId }: { onClose: () => void, parentId?: string }) {
    const mutation = useAddItem();
    const { data: items } = useItems();

    useEffect(() => {
        const handleKeyDown = (e: KeyboardEvent) => {
            if (e.key === 'Escape') {
                onClose();
            }
        };

        document.addEventListener('keydown', handleKeyDown);
        return () => document.removeEventListener('keydown', handleKeyDown);
    }, [onClose]);

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        const form = e.target as HTMLFormElement;
        if (!form.elements) {
            return;
        }

        const titleElement = form.elements.namedItem('title') as HTMLInputElement;
        const parentElement = form.elements.namedItem('parent_id') as HTMLSelectElement;

        // Use either the provided parentId or the selected one from dropdown
        const selectedParentId = parentId || (parentElement?.value !== "none" ? parentElement?.value : undefined);

        const formData = {
            title: titleElement.value,
            parent_id: selectedParentId
        };

        mutation.mutate(formData);
        onClose();
    }

    const modalRef = useRef<HTMLDivElement>(null);

    const handleBackdropClick = (e: React.MouseEvent) => {
        if (modalRef.current && !modalRef.current.contains(e.target as Node)) {
            onClose();
        }
    };

    // Get all potential parent items (all active items)
    const potentialParents = items?.filter(item => !item.completed) || [];

    return (
        <div
            className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50"
            onClick={handleBackdropClick}
        >
            <div
                ref={modalRef}
                className="bg-white rounded-lg shadow-xl p-6 w-full max-w-md"
            >
                <div className="flex justify-between items-center mb-4">
                    <h3 className="text-xl font-semibold">
                        {parentId ? "Add Subtask" : "Add New Task"}
                    </h3>
                    <button
                        onClick={onClose}
                        className="text-gray-500 hover:text-gray-700"
                    >
                        ✕
                    </button>
                </div>

                <form onSubmit={handleSubmit}>
                    <div className="mb-4">
                        <label
                            className="block text-gray-700 text-sm font-bold mb-2"
                            htmlFor="title"
                        >
                            Title
                        </label>
                        <input
                            className="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline"
                            id="title"
                            type="text"
                            placeholder="Enter task title"
                            required
                            autoFocus
                        />
                    </div>

                    {/* Only show parent selection if not creating a subtask */}
                    {!parentId && potentialParents.length > 0 && (
                        <div className="mb-4">
                            <label
                                className="block text-gray-700 text-sm font-bold mb-2"
                                htmlFor="parent_id"
                            >
                                Parent Task (optional)
                            </label>
                            <select
                                className="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline"
                                id="parent_id"
                                defaultValue="none"
                            >
                                <option value="none">None (Top-level task)</option>
                                {potentialParents.map(parent => (
                                    <option key={parent.id} value={parent.id}>
                                        {parent.title}
                                    </option>
                                ))}
                            </select>
                        </div>
                    )}

                    <div className="flex justify-end space-x-2">
                        <button
                            type="button"
                            onClick={onClose}
                            className="bg-gray-300 hover:bg-gray-400 text-gray-800 font-bold py-2 px-4 rounded focus:outline-none focus:shadow-outline"
                        >
                            Cancel
                        </button>
                        <button
                            type="submit"
                            className="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded focus:outline-none focus:shadow-outline"
                        >
                            Create
                        </button>
                    </div>
                </form>
            </div>
        </div>
    );
}
