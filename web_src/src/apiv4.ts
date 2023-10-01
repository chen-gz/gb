import {blogBackendUrl} from "@/config";
import {ref} from "vue";


export interface V4PostData {
    id: number
    title: string
    author: string
    author_email: string
    url: string
    is_draft: boolean
    is_deleted: boolean
    content: string
    summary: string
    tags: string
    category: string
    cover_image: string
    created_at: Date
    updated_at: Date
    view_groups: string[]
    edit_groups: string[]
}


export interface GetPostResponseV3 {
    status: string
    message: string
    post: V4PostData
    html: string
}


export async function getPostV4(url: string, rendered: boolean): Promise<GetPostResponseV3> {
    const request = {
        url: url,
        rendered: rendered
    }
    return await fetch(`${blogBackendUrl}/api/v4/get_post`, {
        method: "POST",
        headers: {
            "Authorization": `Bearer ${localStorage.getItem("token") || ""}`
        },
        body: JSON.stringify(request),
    }).then(response => response.json())
}



export interface UpdatePostResponseV4 {
    status: string
    message: string
    url: string
}

export async function updatePostV4(request: V4PostData): Promise<UpdatePostResponseV4> {
    console.log(localStorage.getItem("token"))
    return await fetch(`${blogBackendUrl}/api/v4/update_post`, {
        method: "POST",
        headers: {
            "Authorization": `Bearer ${localStorage.getItem("token") || ""}`
        },
        body: JSON.stringify(request),
    }).then(response => response.json())
}

export interface NewPostResponseV4 {
    status: string
    message: string
    url: string
}

export async function newPostV4(): Promise<NewPostResponseV4> {
    return await fetch(`${blogBackendUrl}/api/v4/new_post`, {
        method: "POST",
        headers: {
            "Authorization": `Bearer ${localStorage.getItem("token") || ""}`
        },
    }).then(response => response.json())
}

export interface GetDistinctResponse {
    status: string
    message: string
    values: string[]
    length: number
}

export interface GetDistinctRequest {
    field: string
}

export async function getDistinct(col: string): Promise<GetDistinctResponse> {
    return await fetch(`${blogBackendUrl}/api/v4/get_distinct`, {
        method: "POST",
        headers: {
            "Authorization": `Bearer ${localStorage.getItem("token") || ""}`
        },
        body: JSON.stringify({field: col}),
    }).then(response => response.json())
}

export interface SearchPostsRequestV4 {
    author: string
    title: string
    limit: {
        start: number
        size: number
    }
    sort: string
    rendered: boolean
    counts_only: boolean
    content: string
    tags: string
    categories: string
    private_level: number
    is_draft: boolean
    is_deleted: boolean
}

export interface SearchPostsResponseV4 {
    status: string
    message: string
    posts: V4PostData[]
    number_of_posts: number
}

export async function searchPostsV4(request: SearchPostsRequestV4): Promise<SearchPostsResponseV4> {
    return await fetch(`${blogBackendUrl}/api/v4/search_posts`, {
        method: "POST",
        headers: {
            "Authorization": `Bearer ${localStorage.getItem("token") || ""}`
        },
        body: JSON.stringify(request),
    }).then(response => response.json())
}

export interface LoginResponse {
    status: string
    message: string
    email: string
    token: string
    name: string
}

export async function loginV4(email: string, password: string): Promise<LoginResponse> {
    return await fetch(`${blogBackendUrl}/api/v4/login`, {
        method: "POST",
        body: JSON.stringify({
            email: email,
            password: password
        }),
    }).then(response => response.json())
}

// export async function loginV3(){
//     window.location.href= "https://gitea.ggeta.com/login/oauth/authorize?client_id=4093feeb-ff9b-4103-a091-db2381588ce9&redirect_uri=https://blog.ggeta.com&response_type=code&state=STATE"
// }

export async function verifyToken(): Promise<LoginResponse> {
    return await fetch(`${blogBackendUrl}/api/v4/login`, {
        method: "POST",
        headers: {
            "Authorization": `Bearer ${localStorage.getItem("token") || ""}`
        },
    }).then(response => response.json())
}

export let logined = ref(false)

export function logout() {
    localStorage.removeItem("token")
    localStorage.removeItem("userName")
    localStorage.removeItem("userEmail")
    logined.value = false
}

export const alert = ref({
    show: false,
    type: '',
    message: '',
    color: '',
})

export function showError(msg: string) {
    alert.value = {
        show: true,
        type: 'error',
        message: msg,
        color: 'error'
    }
}

export function showSuccess(msg: string) {
    alert.value = {
        show: true,
        type: 'info',
        message: msg,
        color: 'success'
    }
}

export function formatDate(date: Date) {
    date = new Date(date)
    const day = date.toLocaleString("en-US", {day: '2-digit'})
    const month = date.toLocaleString("en-US", {month: 'short'})
    return day + ' ' + month + ' ' + date.getFullYear();
}

