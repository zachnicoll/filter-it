import React, { createRef, useContext, useState } from "react";
import { DropzoneRef, FileWithPath, useDropzone } from "react-dropzone";
import { Filter, PreviewFileWithPath } from "../../api/types";
import API from "../../api";
import { toastError } from "../../common/toast";

import {
  Container,
  FormColumn,
  Img,
  UploadWrapper,
  FormLabel,
  FormArea,
  FormBox,
  FormButton,
  FormButtonWrapper,
  FormCheckBoxInput,
  FormRow,
} from "./UploadBox.styles";
import { filterStringMap } from "../../common/enum.maps";
import { capitalize } from "@material-ui/core";
import { ProgressContext } from "../../hooks/useProgress";

export const UploadBox = () => {
  const { dispatch, state } = useContext(ProgressContext);
  const dropzoneRef = createRef<DropzoneRef>();
  const [file, setFile] = useState<PreviewFileWithPath>();

  const [author, setAuthor] = useState<string | undefined>(undefined);
  const [title, setTitle] = useState<string | undefined>(undefined);
  const [filter, setFilter] = useState<Filter | undefined>(undefined);

  const { getRootProps, getInputProps } = useDropzone({
    maxFiles: 1,
    accept: [".jpg", ".jpeg"],
    // @ts-ignore
    onDrop: (prop: FileWithPath[]) => {
      if (prop[0]) {
        setFile({
          preview: URL.createObjectURL(prop[0]),
          file: prop[0],
        });
      }
    },
  });

  const UploadImage = () => {
    if (file && author && title && filter) {
      API.upload
        .get()
        .then((uploadResponse) => {
          const url = decodeURIComponent(uploadResponse.url);
          API.s3
            .put(url, file.file)
            .then((s3Response) => {
              if (s3Response == 200) {
                API.queue
                  .post({
                    author,
                    title,
                    image: uploadResponse.image,
                    filter: filter,
                  })
                  .then((response) => {
                    dispatch({
                      type: "UPDATE_IMAGE",
                      payload: response.documentID,
                    });
                  })
                  .catch((error) => {
                    toastError(
                      "Image Upload Failed. Please try again later.",
                      error
                    );
                  });
              } else {
                toastError(
                  "Image Upload Failed. Please try again later.",
                  new Error("S3 failed to upload.")
                );
              }
            })
            .catch((error) => {
              toastError("Image Upload Failed. Please try again later.", error);
            });
        })
        .catch((error) => {
          toastError("Image Upload Failed. Please try again later.", error);
        });
    } else {
      toastError("All Fields Required!", new Error("Empty Upload Fields"));
    }
  };

  const openDialog = () => {
    if (dropzoneRef.current) {
      dropzoneRef.current.open();
    }
  };

  return (
    <UploadWrapper>
      {!file ? (
        <Container {...getRootProps()}>
          <input {...getInputProps()} />
          <p>
            Drag &apos;n&apos; drop some files here, or click to select files
          </p>
          <br />
          <button type="button" onClick={openDialog}>
            Open File Dialog
          </button>
        </Container>
      ) : (
        <>
          <Container {...getRootProps()}>
            <input {...getInputProps()} />
            <div onClick={openDialog}>
              <Img src={file.preview} />
            </div>
          </Container>
          <FormArea>
            <FormColumn>
              <FormBox>
                <FormLabel emphasized>Title</FormLabel>
                <input
                  value={title}
                  onChange={(event) => {
                    event.preventDefault();
                    setTitle(event.target.value);
                  }}
                />

                <FormLabel emphasized>Author</FormLabel>
                <input
                  value={author}
                  onChange={(event) => {
                    event.preventDefault();
                    setAuthor(event.target.value);
                  }}
                />
              </FormBox>
            </FormColumn>
            <FormColumn>
              <FormLabel emphasized>Filters</FormLabel>
              <FormBox>
                <FormRow>
                  <FormCheckBoxInput
                    type="radio"
                    name="filter"
                    onChange={() => setFilter(Filter.GRAYSCALE)}
                  />
                  <FormLabel>Grayscale</FormLabel>
                </FormRow>
                <FormRow>
                  <FormCheckBoxInput
                    type="radio"
                    name="filter"
                    onChange={() => setFilter(Filter.SEPIA)}
                  />
                  <FormLabel>Sepia</FormLabel>
                </FormRow>
                <FormRow>
                  <FormCheckBoxInput
                    type="radio"
                    name="filter"
                    onChange={() => setFilter(Filter.INVERT)}
                  />
                  <FormLabel>Invert</FormLabel>
                </FormRow>
              </FormBox>
            </FormColumn>
            <FormButtonWrapper>
              <FormButton onClick={UploadImage}>Upload Image</FormButton>
            </FormButtonWrapper>
          </FormArea>
        </>
      )}
    </UploadWrapper>
  );
};
