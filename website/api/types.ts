import { FileWithPath } from "react-dropzone";

export enum Filter {
  GRAYSCALE,
  SEPIA,
  INVERT,
}

export enum Progress {
  READY,
  PROCESSING,
  DONE,
  FAILED,
}

export interface ImageDocument {
  id: string; // UUID to identify document
  title: string; // Title of art piece
  author: string; // Author of art piece
  filter: Filter; // List of filters to be applied/have been applied to the image
  progress: Progress; // Current progress of image processing in SQS
  image: string; // S3 filename of image related to this document
}

export interface UploadResponse {
  image: string; // S3 filename of image
  url: string; // S3 pre-sign URL
}

export interface PreviewFileWithPath {
  preview: string;
  file: FileWithPath;
}
