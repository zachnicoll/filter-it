import React, { createRef, useState } from "react";
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

export const UploadBox = () => {
  const dropzoneRef = createRef<DropzoneRef>();
  const [file, setFile] = useState<PreviewFileWithPath>();

  const [author, setAuthor] = useState("");
  const [title, setTitle] = useState("");
  const [filter, setFilter] = useState(
    new Array(Object.keys(Filter).length / 2).fill(false)
  );

  const handleChange = (position: number) => {
    const updatedCheckedState = filter.map((item, index) =>
      index === position ? !item : item
    );

    setFilter(updatedCheckedState);
  };

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
    if (file && author && title && !filter.every((v) => v === false)) {
      API.upload.get().then((uploadResponse) => {
        API.s3.put(uploadResponse.url, file.file).then((s3Response) => {
          if (s3Response == 200) {
            API.queue
              .post({
                author,
                title,
                image: uploadResponse.image,
                filter: filter,
              })
              .then((response) => {
                // TODO: push a context to handle Image Process Tracking
              });
          } else {
            toastError(
              "Image Upload Failed. Please try again later.",
              new Error("S3 failed to upload.")
            );
          }
        });
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
                {Object.keys(Filter)
                  // @ts-ignore
                  .filter((key) => !isNaN(Number(Filter[key])))
                  .map((key, index) => (
                    <FormRow key={key}>
                      <FormCheckBoxInput
                        type="checkbox"
                        name={key}
                        value={key}
                        checked={filter[index]}
                        onChange={() => handleChange(index)}
                      />
                      <FormLabel>{key.toLowerCase()}</FormLabel>
                    </FormRow>
                  ))}
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
