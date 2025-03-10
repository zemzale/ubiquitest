import Head from "next/head";
import Link from "next/link";
import { useEffect, useRef, useState } from "react";
import { useAddItem, useCompleteItem, useItems } from "~/query/item";
import { User, useUser } from "~/query/user";
import { useCreateWebsocket, WebSocketProvider } from "~/ws/hook";

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
                <main className="flex min-h-screen flex-col items-center pt-16">
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

    if (!query.data) {
        return <></>;
    }

    const handleComplete = (id: string) => {
        completeMutation.mutate(id);
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

    return <>
        <div className="w-full max-w-md mb-6 mt-6">
            <h2 className="text-xl font-semibold mb-4">Your Items</h2>
            <ul className="space-y-3">
                {activeItems.map((item) => (
                    <li key={item.id}>
                        <div className="bg-white shadow-md rounded-lg p-4 hover:shadow-lg transition-shadow duration-200">
                            <div className="flex items-center justify-between">
                                <div className="flex items-center">
                                    <div>
                                        <h3 className="font-medium text-gray-800">{item.title}</h3>
                                    </div>
                                </div>
                                <button
                                    onClick={() => handleComplete(item.id)}
                                    className="ml-2 bg-green-500 hover:bg-green-600 text-white px-2 py-1 rounded-md text-sm"
                                >
                                    Complete
                                </button>
                            </div>
                        </div>
                    </li>
                ))}

                {hasCompletedItems && (
                    <>
                        <li className="my-6">
                            <div className="relative">
                                <div className="absolute inset-0 flex items-center">
                                    <div className="w-full border-t border-gray-300"></div>
                                </div>
                                <div className="relative flex justify-center">
                                    <span className="bg-white px-2 text-sm text-gray-500">Completed Tasks</span>
                                </div>
                            </div>
                        </li>

                        {completedItems.map((item) => (
                            <li key={item.id}>
                                <div className="bg-white shadow-md rounded-lg p-4 hover:shadow-lg transition-shadow duration-200 opacity-60">
                                    <div className="flex items-center justify-between">
                                        <div className="flex items-center">
                                            <div>
                                                <h3 className="font-medium text-gray-500 line-through">{item.title}</h3>
                                            </div>
                                        </div>
                                        <span className="ml-2 text-green-600 text-sm font-medium">
                                            ✓ Completed
                                        </span>
                                    </div>
                                </div>
                            </li>
                        ))}
                    </>
                )}
            </ul>
        </div>
    </>;
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
    return (
        <nav className="fixed top-0 left-0 right-0 bg-white shadow-md p-4 z-10">
            <div className="container mx-auto flex justify-between items-center">
                <div className="font-bold text-xl">Ubiquitodo</div>
                <div className="flex items-center">
                    <button
                        onClick={onAddTask}
                        className="bg-blue-500 hover:bg-blue-700 text-white font-bold py-1 px-3 rounded text-sm mr-4"
                    >
                        Add Task
                    </button>
                    <p className="mr-4">Logged in as: <span className="font-semibold">{user.username}</span></p>
                    <Link
                        href="/"
                        className="bg-gray-200 hover:bg-gray-300 text-gray-800 font-bold py-1 px-3 rounded text-sm"
                    >
                        Logout
                    </Link>
                </div>
            </div>
        </nav>
    );
}


function AddItemModal({ onClose }: { onClose: () => void }) {
    const mutation = useAddItem();

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
        if (!form.elements.namedItem('title')) {
            return;
        }
        const titleElement = form.elements.namedItem('title') as HTMLInputElement;
        const formData = { title: titleElement.value };

        mutation.mutate(formData);
        onClose();
    }

    const modalRef = useRef<HTMLDivElement>(null);

    const handleBackdropClick = (e: React.MouseEvent) => {
        if (modalRef.current && !modalRef.current.contains(e.target as Node)) {
            onClose();
        }
    };

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
                    <h3 className="text-xl font-semibold">Add New Task</h3>
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
