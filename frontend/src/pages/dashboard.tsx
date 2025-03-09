import Head from "next/head";
import Link from "next/link";
import { useRouter } from "next/router";
import React from "react";
import { useQueryClient } from "@tanstack/react-query";
import { useAddItem, useItems } from "~/query/item";
import { User, useUser } from "~/query/user";
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
    return (
        <>
            <WebSocketProvider value={ws}>
                <main className="flex min-h-screen flex-col items-center justify-center">
                    <h1 className="text-6xl font-bold mb-8">Welcome to Ubiquitodo!</h1>
                    <p className="text-2xl mb-8">
                        You have successfully logged in. {user.username}
                    </p>
                    <ListItems />
                    <AddItem />
                    <Link
                        href="/"
                        className="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded"
                    >
                        Back to Login
                    </Link>
                </main>
            </WebSocketProvider>
        </>
    );
}


function Loading() {
    return (
        <>
            <main className="flex min-h-screen flex-col items-center justify-center">
                <p className="text-2xl mb-8">
                    Loading...
                </p>
            </main>
        </>
    );
}

function ListItems() {
    const query = useItems();
    if (!query.data) {
        return <></>;
    }

    return <>
        <div className="w-full max-w-md mb-6">
            <h2 className="text-xl font-semibold mb-4">Your Items</h2>
            <ul className="space-y-3">
                {query.data.map((item) => (
                    <li key={item.id}>
                        <div className="bg-white shadow-md rounded-lg p-4 hover:shadow-lg transition-shadow duration-200">
                            <div className="flex items-center justify-between">
                                <div className="flex items-center">
                                    <div>
                                        <h3 className="font-medium text-gray-800">{item.title}</h3>
                                    </div>
                                </div>
                            </div>
                        </div>
                    </li>
                ))}
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


function AddItem() {
    const mutation = useAddItem();

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


        titleElement.value = "";
    }

    return <>
        <form
            onSubmit={handleSubmit}
            className="bg-white shadow-md rounded px-8 pt-6 pb-8 mb-4"
        >
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
                    placeholder="Enter your title"
                    required
                />
            </div>
            <button className="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded focus:outline-none focus:shadow-outline w-full">Create</button>
        </form>
    </>;

}
