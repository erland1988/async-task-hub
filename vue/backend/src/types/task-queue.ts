import {Task} from "@/types/task";

export interface TaskQueue {
    id: number;
    app_id: number;
    task_id: number;
    parameters: string;
    relative_delay_time: number;
    delay_execution_time: number;
    execution_status: string;
    execution_status_string: string;
    execution_start: string;
    execution_end: string;
    execution_duration: number;
    execution_count: number;
    created_at: string;
    updated_at: string;
    taskname: string;
    executor_url: string;
    appname: string;
}