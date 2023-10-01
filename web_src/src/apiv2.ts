// import {blogBackendUrl} from "@/config";
// import {ref} from "vue";
//
// export interface PostDataV3Meta {
//     id: number
//     title: string
//     author: string
//     url: string
//     create_time: Date
//     update_time: Date
//     private_level: number
//     summary: string
//     visible_groups: string
//     is_draft: boolean
//     is_deleted: boolean
//     tags: string
//     category: string
//     cover_img: string
// }
//
// export interface PostDataV3Content {
//     id: number
//     content: string
//     category: string
//     tags: string
//
// }
//
// export interface PostDataV3Comment {
//     id: number
//     likes: number
//     dislikes: number
//     viewCount: number
//     comments: string
// }
//
// export interface BlogDataV3 {
//     meta: PostDataV3Meta,
//     content: PostDataV3Content,
//     comment: PostDataV3Comment
// }
//
//
// export interface GetPostResponseV3 {
//     status: string
//     message: string
//     post: BlogDataV3
//     html: string
// }
//
//
// export interface UpdatePostParams {
//     id: number
//     meta: PostDataV3Meta
//     meta_update: boolean
//     content: PostDataV3Content
//     content_update: boolean
//     comment: PostDataV3Comment
//     comment_update: boolean
// }
//
// export interface UpdatePostResponseV3 {
//     status: string
//     message: string
//     url: string
// }
//
// export async function updatePostV3(request: UpdatePostParams): Promise<UpdatePostResponseV3> {
//     console.log(localStorage.getItem("token"))
//     return await fetch(`${blogBackendUrl}/api/v4/update_post`, {
//         method: "POST",
//         headers: {
//             "Authorization": `Bearer ${localStorage.getItem("token") || ""}`
//         },
//         body: JSON.stringify(request),
//     }).then(response => response.json())
// }
//


