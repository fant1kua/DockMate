declare global {
    interface Window {
        runtime?: {
            EventsOn: (event: string, callback: (data: any) => void) => void;
        }
    }
}

export {}; 