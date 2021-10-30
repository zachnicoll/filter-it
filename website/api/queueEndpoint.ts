import axiosInstance from "common/axiosInstance";
import { Filter } from "./types";

interface QueueParams {
  title: string;
  author: string;
  image: string;
  filter: Filter;
}

interface QueueResponse {
  documentID: string;
}

export const post = async (props: QueueParams): Promise<QueueResponse> => {
  const { title, author, filter, image } = props;

  const res = await axiosInstance.post<QueueResponse>(
    "/queue",
    JSON.stringify({
      title: title,
      author: author,
      filter: filter,
      image: image,
    })
  );
  return res.data;
};
