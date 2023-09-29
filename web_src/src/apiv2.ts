import {blogBackendUrl} from "@/config";
import {ref} from "vue";

export interface PostDataV3Meta {
    id: number
    title: string
    author: string
    url: string
    create_time: Date
    update_time: Date
    private_level: number
    summary: string
    visible_groups: string
    is_draft: boolean
    is_deleted: boolean
    tags: string
    category: string
    cover_img: string
}

export interface PostDataV3Content {
    id: number
    content: string
    category: string
    tags: string

}

export interface PostDataV3Comment {
    id: number
    likes: number
    dislikes: number
    viewCount: number
    comments: string
}

export interface BlogDataV3 {
    meta: PostDataV3Meta,
    content: PostDataV3Content,
    comment: PostDataV3Comment
}


export interface GetPostResponseV3 {
    status: string
    message: string
    post: BlogDataV3
    html: string
}

export async function getPostV3(url: string, rendered: boolean): Promise<GetPostResponseV3> {
    const request = {
        url: url,
        rendered: rendered
    }
    return await fetch(`${blogBackendUrl}/api/v3/get_post`, {
        method: "POST",
        headers: {
            "Authorization": `Bearer ${localStorage.getItem("token") || ""}`
        },
        body: JSON.stringify(request),
    }).then(response => response.json())
}


export interface UpdatePostParams {
    id: number
    meta: PostDataV3Meta
    meta_update: boolean
    content: PostDataV3Content
    content_update: boolean
    comment: PostDataV3Comment
    comment_update: boolean
}

export interface UpdatePostResponseV3 {
    status: string
    message: string
    url: string
}

export async function updatePostV3(request: UpdatePostParams): Promise<UpdatePostResponseV3> {
    console.log(localStorage.getItem("token"))
    return await fetch(`${blogBackendUrl}/api/v3/update_post`, {
        method: "POST",
        headers: {
            "Authorization": `Bearer ${localStorage.getItem("token") || ""}`
        },
        body: JSON.stringify(request),
    }).then(response => response.json())
}

export interface NewPostResponseV3 {
    status: string
    message: string
    url: string
}

export async function newPostV3(): Promise<NewPostResponseV3> {
    return await fetch(`${blogBackendUrl}/api/v3/new_post`, {
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
    column: string
}

export async function getDistinct(col: string): Promise<GetDistinctResponse> {
    return await fetch(`${blogBackendUrl}/api/v3/get_distinct`, {
        method: "POST",
        headers: {
            "Authorization": `Bearer ${localStorage.getItem("token") || ""}`
        },
        body: JSON.stringify({column: col}),
    }).then(response => response.json())
}

export interface SearchPostsRequestV3 {
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

export interface SearchPostsResponseV3 {
    status: string
    message: string
    posts: PostDataV3Meta[]
    number_of_posts: number
}

export async function searchPostsV3(request: SearchPostsRequestV3): Promise<SearchPostsResponseV3> {
    return await fetch(`${blogBackendUrl}/api/v3/search_posts`, {
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

export async function loginV3(email: string, password: string): Promise<LoginResponse> {
    return await fetch(`${blogBackendUrl}/api/v3/login`, {
        method: "POST",
        headers: {
            "Authorization": `Basic ${email}:${password}`
        },
    }).then(response => response.json())
}

// export async function loginV3(){
//     window.location.href= "https://gitea.ggeta.com/login/oauth/authorize?client_id=4093feeb-ff9b-4103-a091-db2381588ce9&redirect_uri=https://blog.ggeta.com&response_type=code&state=STATE"
// }

export async function verifyToken(): Promise<LoginResponse> {
    return await fetch(`${blogBackendUrl}/api/v3/login`, {
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

