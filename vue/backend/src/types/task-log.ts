
export interface TaskLog {
    id: number;
    app_id: number;
    task_id: number;
    task_queue_id: number;
    request_id: string;
    action: string;
    message: string;
    created_at: string;
}