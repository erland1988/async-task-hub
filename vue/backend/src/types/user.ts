
export interface User {
    id: number;
    username: string;
    password: string;
    truename: string;
    role: string;
    expires_at: string;
    created_at: string;
    updated_at: string;
}

export interface Register {
    username: string;
    password: string;
}