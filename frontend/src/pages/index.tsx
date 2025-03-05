import Head from "next/head";
import { useRouter } from "next/router";
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
    const mutation = useLogin()
    const router = useRouter();

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();

        const form = e.target;
        const formData = {
            username: form.elements.username.value,
        };

        mutation.mutate(formData);

    }


    if (mutation.isSuccess) {
        router.push("/dashboard");
        return <></>
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
                />
                <button className="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded focus:outline-none focus:shadow-outline w-full">
                    Login
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
