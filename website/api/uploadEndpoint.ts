import axiosInstance from "common/axiosInstance";
import { UploadResponse } from "./types";

export const get = async (): Promise<UploadResponse> => {
    const res = await axiosInstance.get<UploadResponse>("/upload");
    return res.data;
};
