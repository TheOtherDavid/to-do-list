export interface TaskInstance {
    id: string;
    template_id: string;
    title: string;
    description: string;
    completed: boolean;
    created_at: string; // can be Date later, but JSON will send string
    completed_at: string | null;
  }
  