import {showSuccess} from "@/apiv4";

const photoBackendUrl = "https://blog.ggeta.com"

import * as crypto from "crypto";
import {compileString} from "sass";

// const photoBackendUrl = "http://localhost:2009"

export interface PhotoItemV2 {
    id: number,
    ori_hash: string,
    jpg_md5: string,
    jpg_sha256: string,
    thumb_hash: string,
    has_original: boolean,
    ori_ext: string,
    deleted: boolean,
    tags: string,
    category: string,
}

// export interface GetPhotoRequest {
//     id: number,
// }

export interface GetPhotoResponse {
    photo: PhotoItemV2,
    thum_url: string,
    ori_url: string,
    jpeg_url: string,
    message: string,
}

var cnt = 0

export async function getPhotoById(id: number): Promise<GetPhotoResponse> {
    return await fetch(`${photoBackendUrl}/api/photo/v2/get_photo_id/` + id, {
        method: "GET",
        headers: {
            "Authorization": `Bearer ${localStorage.getItem("token") || ""}`
        },
    }).then(response => {
        console.log("getPhoto", cnt)
        cnt += 1
        return response.json()
    })
}

export async function getPhotoByHash(hash: string): Promise<GetPhotoResponse> {
    return await fetch(`${photoBackendUrl}/api/photo/v2/get_photo_hash/` + hash, {
        method: "GET",
        headers: {
            "Authorization": `Bearer ${localStorage.getItem("token") || ""}`
        },
    }).then(response => {
        console.log("getPhoto", cnt)
        cnt += 1
        return response.json()
    })
}


export interface InsertPhotoResponse {
    id: number,
    message: string,
    presigned_original_url: string,
    presigned_thumb_url: string,
    presigned_jpeg_url: string,
}

export async function insertPhoto(request: PhotoItemV2): Promise<InsertPhotoResponse> {
    return await fetch(`${photoBackendUrl}/api/photo/v2/insert_photo`, {
        method: "POST",
        headers: {
            "Authorization": `Bearer ${localStorage.getItem("token") || ""}`
        },
        body: JSON.stringify(request),
    }).then(response => {
        return response.json()
    })
}

export async function updatePhotoFile(request: PhotoItemV2): Promise<InsertPhotoResponse> {
    return await fetch(`${photoBackendUrl}/api/photo/v2/update_photo_file`, {
        method: "POST",
        headers: {
            "Authorization": `Bearer ${localStorage.getItem("token") || ""}`
        },
        body: JSON.stringify(request),
    }).then(response => {
        return response.json()
    })
}

function arrayBufferToHexString(buffer: ArrayBuffer) {
    const uint8Array = new Uint8Array(buffer);
    return Array.from(uint8Array, byte => byte.toString(16).padStart(2, '0')).join('');
}

async function calculateMD5(input: Uint8Array): Promise<string> {
    // @ts-ignore
    return hashwasm.md5(input)
}

async function calculateSHA256(input: Uint8Array): Promise<string> {
    // @ts-ignore
    return hashwasm.sha256(input)
}


// ori_file can be nil or File
export async function addPhoto(jpeg_file: File, ori_file?: File) {
    const request = {} as PhotoItemV2
    // resize jpeg file
    const thumbnail_file = await resizeImage(jpeg_file)
    // calculate hash of thumbnail file
    // request.thumb_hash
    var thumb_hash = thumbnail_file.arrayBuffer().then((buffer) => {
        const uint8array = new Uint8Array(buffer);
        return calculateMD5(uint8array);
    })
    // request.jpg_md5
    var jpg_md5 = jpeg_file.arrayBuffer().then((buffer) => {
        const uint8array = new Uint8Array(buffer);
        return calculateMD5(uint8array);
    })
    // request.jpg_sha256
    var jpgsha256 = jpeg_file.arrayBuffer().then((buffer) => {
        const uint8array = new Uint8Array(buffer);
        return calculateSHA256(uint8array);
    })
    // wait for all hash calculation
    await Promise.all([thumb_hash, jpg_md5, jpgsha256]).then(
        (values) => {
            request.thumb_hash = values[0]
            request.jpg_md5 = values[1]
            request.jpg_sha256 = values[2]
        }
    )
    // calculate hash of orifile if exists
    // request.ori_hash
    if (ori_file != null) {
        request.has_original = true
        request.ori_ext = ori_file.name.split(".")[1]
        request.ori_hash = await ori_file.arrayBuffer().then((buffer) => {
            const uint8array = new Uint8Array(buffer);
            return calculateMD5(uint8array);
        })
    }
    else {
        request.has_original = false
        request.ori_ext = ""
        request.ori_hash = ""
    }

    const update_response = await insertPhoto(request)
    console.log(request)
    console.log(update_response)
    // the first two characters of the hash is ok
    if (update_response.message.substring(0, 2) == "ok") {
        var ori_upload_promise: Promise<Response> = Promise.resolve(new Response())
        var jpeg_upload_promise: Promise<Response> = Promise.resolve(new Response())
        var thumb_upload_promise: Promise<Response> =  Promise.resolve(new Response())
        if (update_response.presigned_original_url.length > 0 && ori_file != null) {
            ori_upload_promise = uploadFileToPresignURL(ori_file, update_response.presigned_original_url)
        }
        if (update_response.presigned_jpeg_url.length > 0) {
            jpeg_upload_promise = uploadFileToPresignURL(jpeg_file, update_response.presigned_jpeg_url)
        }
        if (update_response.presigned_thumb_url.length > 0) {
            thumb_upload_promise = uploadFileToPresignURL(thumbnail_file, update_response.presigned_thumb_url)
        }
        await Promise.all([ori_upload_promise, jpeg_upload_promise, thumb_upload_promise])
    }
}

function uploadFileToPresignURL(file: File, presignedURL: string): Promise<Response> {
    return fetch(presignedURL, {
        method: 'PUT',
        body: file
    })
}

async function resizeImage(inputFile: File): Promise<File> {
    return new Promise<File>((resolve) => {
        const img = new Image();
        img.onload = () => {
            const canvas = document.createElement('canvas');
            const ctx = canvas.getContext('2d');
            canvas.height = 256;
            // set width to 256 * ratio
            canvas.width = 256 * img.width / img.height;

            // @ts-ignore
            ctx.drawImage(img, 0, 0, canvas.width, canvas.height);

            canvas.toBlob((blob) => {
                // @ts-ignore
                const resizedFile = new File([blob], inputFile.name, {
                    type: 'image/jpeg',
                });

                resolve(resizedFile);
            }, 'image/jpeg', 1); // Quality 1 means no compression (optional)
        };

        const reader = new FileReader();
        reader.onload = (event) => {
            img.src = event.target?.result as string;
        };
        reader.readAsDataURL(inputFile);
    });
}


export async function uploadPhotos(files: FileList) {
    var ori_files = new Map()
    var jpeg_files = new Map()
    var heic_files = new Map()
    for (var i = 0; i < files.length; i++) {
        var file = files[i];
        var filename_without_ext = file.name.split(".")[0]
        if (file.name.toLowerCase().endsWith(".nef")) {
            ori_files.set(filename_without_ext, file)
        } else if (file.name.toLowerCase().endsWith(".heic")) { // iPhone
            heic_files.set(filename_without_ext, file)
        } else if (file.name.toLowerCase().endsWith(".jpg")
            || file.name.endsWith(".jpeg")
            // || file.name.endsWith(".png")
        ) {
            jpeg_files.set(filename_without_ext, file)
        }
    }
    // only original file is not accept
    for (var [key, value] of ori_files) {
        if (jpeg_files.has(key) == false) {
            alert("Please upload the original file and the corresponding jpeg file together. The server does not process the original file alone.")
            return
        }
    }
    showSuccess("Total " + jpeg_files.size + " photos to upload. Uploading...")
    // upload file to server
    let cnt = 0
    const tasks = []
    for (var [key, value] of jpeg_files) {
        var ori_file = ori_files.get(key)
        tasks.push(addPhoto(value, ori_file))
        // await addPhoto(value, ori_file)
        if (tasks.length >= 10) {
            await Promise.all(tasks)
            // clear tasks
            tasks.length = 0
            showSuccess("Uploaded " + cnt + " photos. " + (jpeg_files.size + heic_files.size - cnt) + " photos left.")
        }
        cnt += 1
    }
    if (tasks.length > 0){
        await Promise.all(tasks)
        showSuccess("Uploaded " + cnt + " photos. " + (jpeg_files.size + heic_files.size - cnt) + " photos left.")
    }
    // upload heic file to server
    // for (var [key, value] of heic_files) {
    //     await addPhoto(value)
    //     // number of photos uploaded and how many photos left
    //     showSuccess("Uploaded " + cnt + " photos. " + (jpeg_files.size + heic_files.size - cnt) + " photos left.")
    //     cnt += 1
    // }
}

export interface PhotoListResponse {
    message: string,
    ids: number[],
}

export async function getPhotoIds(): Promise<PhotoListResponse> {
    return await fetch(`${photoBackendUrl}/api/photo/v2/get_photo_list`, {
        method: "GET",
        headers: {
            "Authorization": `Bearer ${localStorage.getItem("token") || ""}`
        },
    }).then(response => response.json())
}

//
// export async function UpdatePhoto(photo: PhotoItemV2){
//     return await fetch(`${photoBackendUrl}/api/photo/v1/update_photo`, {
//         method: "POST",
//         headers: {
//             "Authorization": `Bearer ${localStorage.getItem("token") || ""}`
//         },
//         body: JSON.stringify(photo),
//     }).then(response => response.json())
// }
//
// export async function getDeletedPhotoIds(): Promise<PhotoListResponse> {
//     return await fetch(`${photoBackendUrl}/api/photo/v1/get_deleted_photo_list`, {
//         method: "GET",
//         headers: {
//             "Authorization": `Bearer ${localStorage.getItem("token") || ""}`
//         },
//     }).then(response => response.json())
// }
