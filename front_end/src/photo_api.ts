const photoBackendUrl = "https://blog.ggeta.com"

import * as crypto from "crypto";
import {compileString} from "sass";

// const photoBackendUrl = "http://localhost:2009"

export interface PhotoItem {
    id: number,
    hash: string,
    has_original: boolean,
    original_ext: string,
    deleted: boolean,
    tags: string,
    category: string,
}

export interface GetPhotoRequest {
    id: number,
}

export interface GetPhotoResponse {
    photo: PhotoItem,
    thum_url: string,
    ori_url: string,
    jpeg_url: string,
    message: string,
}

var cnt = 0
export async function getPhoto(request: GetPhotoRequest): Promise<GetPhotoResponse> {
    return await fetch(`${photoBackendUrl}/api/photo/v1/get_photo`, {
        method: "POST",
        headers: {
            "Authorization": `Bearer ${localStorage.getItem("token") || ""}`
        },
        body: JSON.stringify(request),
    }).then(response => {
        console.log("getPhoto", cnt)
        cnt += 1
        return response.json()
    })
}

export interface InsertPhotoRequest {
    hash: string,
    has_original: boolean,
    original_ext: string,
}

export interface InsertPhotoResponse {
    message: string,
    presigned_original_url: string,
    presigned_thumbnail_url: string,
    presigned_jpeg_url: string,
}

export async function insertPhoto(request: InsertPhotoRequest): Promise<InsertPhotoResponse> {
    return await fetch(`${photoBackendUrl}/api/photo/v1/insert_photo`, {
        method: "POST",
        headers: {
            "Authorization": `Bearer ${localStorage.getItem("token") || ""}`
        },
        body: JSON.stringify(request),
    }).then(response => {
        return response.json()
    })
}

// import * as crypto from 'crypto';
function arrayBufferToHexString(buffer: ArrayBuffer) {
    const uint8Array = new Uint8Array(buffer);
    return Array.from(uint8Array, byte => byte.toString(16).padStart(2, '0')).join('');
}

async function calculateSHA1(input: Uint8Array): Promise<string> {
    // console.log(crypto)
    console.log(input)
    // convert input to arraybuffer
    const sha1Hash = await window.crypto.subtle.digest('sha-1', input)
    // createHash('sha1');
    // return  String.fromCharCode(...new Uint8Array(sha1Hash))
    // sha1Hash.update(Buffer.from(input));
    return arrayBufferToHexString(sha1Hash)

    // return sha1Hash.digest('hex');
}


// ori_file can be nil or File
export async function addPhoto(jpeg_file: File, ori_file?: File) {
    var request = {} as InsertPhotoRequest
    if (ori_file != null) {
        // calculate hash of orifile
        request.hash = await ori_file.arrayBuffer().then((buffer) => {
            const uint8array = new Uint8Array(buffer);
            return calculateSHA1(uint8array);
        })
        request.has_original = true
    } else {
        // calculate hash of jpegfile
        request.hash = await jpeg_file.arrayBuffer().then((buffer) => {
            const uint8array = new Uint8Array(buffer);
            return calculateSHA1(uint8array);
        })
        request.has_original = false
    }
    console.log("request", request)
    // return
    request.original_ext = ori_file?.name.split('.').pop() || ""
    const response = await insertPhoto(request)
    if (response.message != "ok") {
        console.log("response", response)
        console.error(response.message)
        return
    }
    console.log(response)
    if (ori_file != null) {
        await uploadFileToPresignURL(ori_file, response.presigned_original_url)
    }
    await uploadFileToPresignURL(jpeg_file, response.presigned_jpeg_url)
    // process thumbnail from jpeg_file
    const thumbnail_file = await resizeImage(jpeg_file)
    await uploadFileToPresignURL(thumbnail_file, response.presigned_thumbnail_url)
}

function uploadFileToPresignURL(file: File, presignedURL: string): Promise<Response> {
    return fetch(presignedURL, {
        method: 'PUT',
        body: file
    })
}

// follwing code generate by chatgpt.  // todo: verify
async function resizeImage(inputFile: File): Promise<File> {
    return new Promise<File>((resolve) => {
        const img = new Image();
        img.onload = () => {
            const canvas = document.createElement('canvas');
            const ctx = canvas.getContext('2d');
            canvas.width = 256;
            canvas.height = 256;

            // @ts-ignore
            ctx.drawImage(img, 0, 0, 256, 256);

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

    // for Nikon format
    for (var i = 0; i < files.length; i++) {
        var file = files[i];
        var filename_without_ext = file.name.split(".")[0]
        console.log("filename_without_ext: " + filename_without_ext)
        if (file.name.endsWith(".NEF")) {
            ori_files.set(filename_without_ext, file)
        } else if (file.name.endsWith(".JPG")) {
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

    // upload file to server
    for (var [key, value] of jpeg_files) {
        var ori_file = ori_files.get(key)
        await addPhoto(value, ori_file)
    }
}

export interface PhotoListResponse {
    message: string,
    ids: number[],
}

export async function getPhotoIds(): Promise<PhotoListResponse> {
    return await fetch(`${photoBackendUrl}/api/photo/v1/get_photo_list`, {
        method: "GET",
        headers: {
            "Authorization": `Bearer ${localStorage.getItem("token") || ""}`
        },
    }).then(response => response.json())
}

export async function UpdatePhoto(photo: PhotoItem){
    return await fetch(`${photoBackendUrl}/api/photo/v1/update_photo`, {
        method: "POST",
        headers: {
            "Authorization": `Bearer ${localStorage.getItem("token") || ""}`
        },
        body: JSON.stringify(photo),
    }).then(response => response.json())
}

export async function getDeletedPhotoIds(): Promise<PhotoListResponse> {
    return await fetch(`${photoBackendUrl}/api/photo/v1/get_deleted_photo_list`, {
        method: "GET",
        headers: {
            "Authorization": `Bearer ${localStorage.getItem("token") || ""}`
        },
    }).then(response => response.json())
}
