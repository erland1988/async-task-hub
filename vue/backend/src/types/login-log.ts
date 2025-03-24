import {User} from "@/types/user";

export interface LoginLog {
    id: number;
    admin_id: number;
    token: string;
    admin: User;
    expires_at: string;
    created_at: string;
}