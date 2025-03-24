
export interface User {
    id: number;
    username: string;
    password: string;
    phone: string;
    email: string;
    truename: string;
    rolename: string;
    role: string;
    expires_at: string;
    created_at: string;
    updated_at: string;
}

export interface Register {
    username: string;
    password: string;
}