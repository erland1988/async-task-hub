import {User} from "@/types/user";

export interface Log {
    id: number;
    admin_id: number;
    operation: string;
    details: string;
    admin: User;
    created_at: string;
}