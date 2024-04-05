import { useLoaderData } from "react-router-dom";
import { Link } from "react-router-dom";
import * as React from "react";

const api = "http://127.0.0.1:8084/api";
const photoCDN = "https://cdn.lumi.mitsukaki.com/";

let album, user;

export default function UploadPage() {
    const result = useLoaderData();
    album = result.album;
    user = result.user;

    return (
        <div class="flex min-h-full flex-col justify-center px-6 py-12 lg:px-8">
            <div class="sm:mx-auto sm:w-full sm:max-w-sm">
                <img class="mx-auto h-10 w-auto" src="https://tailwindui.com/img/logos/mark.svg?color=indigo&shade=600" alt="Your Company" />
                <h2 class="mt-10 text-center text-2xl font-bold leading-9 tracking-tight text-gray-900">Upload Images</h2>
            </div>
            <div className="mt-10 sm:mx-auto sm:w-full sm:max-w-sm">
                <form className="space-y-6" onSubmit={uploadImages}>
                    <div className="mb-4">
                        <label htmlFor="images" className="block text-lg font-medium text-gray-700">Images</label>
                        <input id="images" type="file" name="files" multiple />
                    </div>
                    <div>
                        <button type="submit" class="flex w-full justify-center rounded-md bg-indigo-600 px-3 py-1.5 text-sm font-semibold leading-6 text-white shadow-sm hover:bg-indigo-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600">Upload</button>
                    </div>
                </form>
            </div>
        </div>
    );
}

async function uploadImages(e) {
    e.preventDefault();

    
    // Get form data
    const files = e.target.elements.files.files;
    
    // Upload every file in the files array
    let uploaded = 0;
    for (let i = 0; i < files.length; i++) {
        const file = files[i];
        const formData = new FormData();
        formData.append("image", file);
        formData.append("title", file.name);
        formData.append("description", "No description");
        formData.append("aspect_ratio", 1.0);
        formData.append("row", i);

        // Upload image
        try {
            const response = await fetch(`${api}/album/${album._id}`, {
                method: "PUT",
                body: formData
            });
            if (response.ok) {
                console.log(file.name + " uploaded successfully");
                uploaded++;

                // Redirect to album page if all images are uploaded
                if (uploaded === files.length) {
                    window.location.href = `/a/${album._id}`;
                }
            } else {
                console.log("Error uploading image");
            }
        } catch (error) {
            console.log(error);
            console.log("Error uploading image");
        }
    }
}

export async function UploadPageLoader({ params }) {
    const albumId = params.album_id;
    const response = await fetch(`${api}/album/${albumId}`);
    const album = await response.json();

    const userId = album.author_user_id;
    const userResponse = await fetch(`${api}/user/${userId}`);
    const user = await userResponse.json();

    return { 
        album: album,
        user: user
    };
}
