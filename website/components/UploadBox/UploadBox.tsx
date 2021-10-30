import React, { createRef, useState } from "react";
import { DropzoneRef, FileWithPath, useDropzone } from "react-dropzone";
import { Filter, PreviewFileWithPath } from "../../api/types";
import API from "../../api";
import { toastError } from "../../common/toast";
import * as Styles from "./UploadBox.styles";
import { filterRadioButtons } from "common/util";
import { DragNDrop } from "./molecules/DragNDrop";
import { useUploadImage } from "hooks/useUploadImage";
import { Spinner } from "components";

export const UploadBox = () => {
  const {upload, uploading} = useUploadImage();

  const dropzoneRef = createRef<DropzoneRef>();
  const [file, setFile] = useState<PreviewFileWithPath>();

  const [author, setAuthor] = useState<string | undefined>(undefined);
  const [title, setTitle] = useState<string | undefined>(undefined);
  const [filter, setFilter] = useState<Filter | undefined>(undefined);

  const handleDrop = (files: FileWithPath[]) => {
    if (files.length === 0) {
      toastError(
        "Only .jpg and .jpeg files are supported!",
        new Error("Unsupported image type")
      );
    }

    if (files[0]) {
      setFile({
        preview: URL.createObjectURL(files[0]),
        file: files[0],
      });
    }
  };

  const { getRootProps, getInputProps, isDragActive } = useDropzone({
    maxFiles: 1,
    accept: [".jpg", ".jpeg"],
    onDrop: handleDrop,
  });

  const uploadImage = async () => {
    if (file && author && title && filter) {
      await upload(author, title, filter, file);
    } else {
      toastError("All Fields Required!", new Error("Empty Upload Fields"));
    }
  };

  const handleOpenDialog = () => {
    if (dropzoneRef.current) {
      dropzoneRef.current.open();
    }
  };

  return (
    <Styles.UploadWrapper>
      {!file ? (
        <DragNDrop
          getInputProps={getInputProps}
          getRootProps={getRootProps}
          openDialog={handleOpenDialog}
          dragActive={isDragActive}
        />
      ) : (
        <>
          <Styles.Container {...getRootProps()}>
            <input {...getInputProps()} />
            <div onClick={handleOpenDialog}>
              {uploading ? <Spinner type="Oval"/> : <Styles.Img src={file.preview} />}
            </div>
          </Styles.Container>

          <Styles.FormArea>
            <Styles.FormRow>
              <Styles.FormBox>
                <Styles.FormLabel emphasized>Title</Styles.FormLabel>

                <Styles.Input
                  value={title}
                  onChange={(event) => {
                    event.preventDefault();
                    setTitle(event.target.value);
                  }}
                  placeholder="My Beautiful Filtered Image"
                />
              </Styles.FormBox>

              <Styles.FormBox>
                <Styles.FormLabel emphasized>Author</Styles.FormLabel>

                <Styles.Input
                  value={author}
                  onChange={(event) => {
                    event.preventDefault();
                    setAuthor(event.target.value);
                  }}
                  placeholder="Zadi Nailey"
                />
              </Styles.FormBox>
            </Styles.FormRow>

            <Styles.FormColumn>
              <Styles.FormLabel emphasized>Select Filter</Styles.FormLabel>

              <Styles.FormRow>
                {filterRadioButtons.map((radioBtnValue) => (
                  <Styles.FormRow key={radioBtnValue.value}>
                    <Styles.FormCheckBoxInput
                      type="radio"
                      name="filter"
                      onChange={() => setFilter(radioBtnValue.value)}
                    />
                    <Styles.FormLabel>{radioBtnValue.label}</Styles.FormLabel>
                  </Styles.FormRow>
                ))}
              </Styles.FormRow>
            </Styles.FormColumn>

            <Styles.FormButton onClick={uploadImage}>
              Upload Image
            </Styles.FormButton>
          </Styles.FormArea>
        </>
      )}
    </Styles.UploadWrapper>
  );
};
