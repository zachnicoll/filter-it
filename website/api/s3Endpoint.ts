import { FileWithPath } from "react-dropzone";
import axios from "axios";

export const put = async (url: string, image: FileWithPath): Promise<void> => {
  await axios.put(url, image);
};
