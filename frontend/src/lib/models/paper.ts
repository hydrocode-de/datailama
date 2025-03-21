export interface Paper {
    id: number;
    title: string;
    doi: string;
    author: string;
    journal: string;
    url: string;
    published: Date;
    citations: number;
    citations_year: number;
}