import Head from "next/head";
import Link from "next/link";
import { useRouter } from "next/router";
import { handleClientScriptLoad } from "next/script";
import { useAddItem } from "~/query/item";
import { useUser } from "~/query/user";

export default function Dashboard() {
    const router = useRouter();
    const query = useUser();

    if (query.isLoading) {
        return <Loading />
    }

    if (query.isError) {
        router.push("/");
        return <></>;
    }

    if (!query.data) {
        router.push("/");
        return <></>;
    }


    return (
        <>
            <Head>
                <title>Dashboard | Ubiquitodo</title>
            </Head>
            <main className="flex min-h-screen flex-col items-center justify-center">
                <h1 className="text-6xl font-bold mb-8">Welcome to Ubiquitodo!</h1>
                <p className="text-2xl mb-8">
                    You have successfully logged in. {query.data.username}
                </p>
                <AddItem />
                <Link
                    href="/"
                    className="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded"
                >
                    Back to Login
                </Link>
            </main>
        </>
    );
}


function Loading() {
    return (
        <>
            <main className="flex min-h-screen flex-col items-center justify-center">
                <h1 className="text-6xl font-bold mb-8">Welcome to Ubiquitodo!</h1>
                <p className="text-2xl mb-8">
                    Loading...
                </p>
            </main>
        </>
    );
}


function AddItem() {
    const mutation = useAddItem();

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        const form = e.target;
        const formData = {
            title: form.elements.title.value,
        };

        mutation.mutate(formData);
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
