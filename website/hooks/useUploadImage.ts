import API from "api";
import { Filter, PreviewFileWithPath, Progress } from "api/types";
import { toastError } from "common/toast";
import { useProgress } from "context";
import { useState } from "react";

interface HookReturn {
  upload: (
    author: string,
    title: string,
    tag: Filter,
    file: PreviewFileWithPath
  ) => Promise<void>;
  uploading: boolean;
  dispatchUploading: (uploading: boolean) => void;
}

export const useUploadImage = (): HookReturn => {
  const { dispatchProgress } = useProgress();
  const [uploading, setUploading] = useState(false);

  const upload = async (
    author: string,
    title: string,
    tag: Filter,
    file: PreviewFileWithPath
  ): Promise<void> => {
    // Retrieve signed URL to upload to S3
    const uploadResponse = await API.upload.get();
    const url = decodeURIComponent(uploadResponse.url);

    // Upload to S3 through signed URL
    await API.s3.put(url, file.file);

    // Create new document and queue up image for processing
    const queueResponse = await API.queue.post({
      author,
      title,
      tag,
      image: uploadResponse.image,
    });

    dispatchProgress({
      type: "UPDATE_IMAGE",
      payload: { id: queueResponse.documentID },
    });
  };

  return { upload, uploading, dispatchUploading: setUploading };
};
