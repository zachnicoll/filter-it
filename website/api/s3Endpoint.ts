import axiosInstance from "../common/axiosInstance";
import { FileWithPath } from "react-dropzone";
import axios from "axios";

export const put = async (
  url: string,
  image: FileWithPath
): Promise<number> => {
  const res = await axios.put(url, image);
  return res.status;
};
