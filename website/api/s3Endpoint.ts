import axiosInstance from "../common/axiosInstance";
import { FileWithPath } from "react-dropzone";

export const put = async (
  url: string,
  image: FileWithPath
): Promise<number> => {
  const res = await axiosInstance.put(url, image);
  return res.status;
};
