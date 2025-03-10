import Head from "next/head";
import { useRouter } from "next/router";
import React, { useEffect, useRef, useState } from "react";
import { useLogin } from "~/query/user";

export default function Home() {
    return (
        <>
            <Head>
                <title>Ubiquitodo</title>
            </Head>
            <main className="flex min-h-screen flex-col items-center justify-center">
                <h1 className="text-6xl font-bold">Ubiquitodo</h1>
                <p className="text-2xl mb-8">
                    A simple todo app built with Ubiquitous
                </p>
                <div className="w-full max-w-md">
                    <LoginFrom />
                </div>
            </main>
        </>
    );
}

function LoginFrom() {
    const mutation = useLogin();
    const router = useRouter();
    const [isRedirecting, setIsRedirecting] = useState(false);

    // Handle successful login and navigation to dashboard
    useEffect(() => {
        // Only run this effect if login was successful and we're not already redirecting
        if (mutation.isSuccess && !isRedirecting) {
            console.log('Login successful, redirecting to dashboard');
            setIsRedirecting(true);
            router.push("/dashboard");
        }
    }, [mutation.isSuccess, router, isRedirecting]);

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();

        // Prevent duplicate submissions
        if (mutation.isPending || isRedirecting) {
            console.log('Preventing duplicate login submission');
            return;
        }


        const form = e.target as HTMLFormElement;
        const usernameInput = form.elements.namedItem('username') as HTMLInputElement;
        const username = usernameInput.value.trim();

        console.log('Submitting login form with username:', username);
        mutation.mutate({ username });
    }

    return <>
        <form
            onSubmit={handleSubmit}
            className="bg-white shadow-md rounded px-8 pt-6 pb-8 mb-4"
        >
            <div className="mb-4">
                <label
                    className="block text-gray-700 text-sm font-bold mb-2"
                    htmlFor="username"
                >
                    Username
                </label>
                <input
                    className="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline"
                    id="username"
                    type="text"
                    placeholder="Enter your username"
                    required
                    disabled={mutation.isPending}
                />
                <button
                    className="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded focus:outline-none focus:shadow-outline w-full"
                    disabled={mutation.isPending}
                >
                    {mutation.isPending ? 'Logging in...' : 'Login'}
                </button>
                {mutation.isError ? <div className="mt-4">
                    <div className="mb-4 text-red-500 text-sm">
                        {mutation.error.message}
                    </div>
                </div> : null}
                {mutation.isPending ? <div className="mt-4">Loading...</div> : null}
            </div>
        </form>
    </>;
}
