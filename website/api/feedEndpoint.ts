import axiosInstance from "common/axiosInstance";
import { Filter, ImageDocument } from "./types";

const ENDPOINT = "/feed";

export const get = async (filter: {
  value: Filter | null;
  label: string;
}): Promise<ImageDocument[]> => {
  const res = await axiosInstance.get<ImageDocument[]>(
    filter?.value ? `${ENDPOINT}?filter=${filter.value}` : ENDPOINT
  );
  return res.data;
};
