import { PageWrapper } from "components";
import type { NextPage } from "next";
import React, { createRef, ReactNode, useState } from "react";
import Dropzone, {
  DropzoneRef,
  FileWithPath,
  useDropzone,
} from "react-dropzone";
import styled from "@emotion/styled";

const thumbsContainer = {
  display: "flex",
  flexDirection: "row",
  flexWrap: "wrap",
  marginTop: 16,
};

const thumb = {
  display: "inline-flex",
  borderRadius: 2,
  border: "1px solid #eaeaea",
  marginBottom: 8,
  marginRight: 8,
  width: 100,
  height: 100,
  padding: 4,
  boxSizing: "border-box",
};

const thumbInner = {
  display: "flex",
  minWidth: 0,
  overflow: "hidden",
};

const img = {
  display: "block",
  width: "auto",
  height: "100%",
};

const Container = styled.div`
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 20px;
  border-width: 2px;
  border-radius: 2px;
  border-color: #eeeeee;
  border-style: dashed;
  background-color: #fafafa;
  color: #bdbdbd;
  outline: none;
  transition: border 0.24s ease-in-out;
`;

const Upload: NextPage = () => {
  const dropzoneRef = createRef<DropzoneRef>();
  const [files, setFiles] = useState<FileWithPath[]>([]);
  const { getRootProps, getInputProps } = useDropzone({
    maxFiles: 1,
    accept: [".jpg", ".jpeg"],
    onDrop: (acceptedFiles: FileWithPath[]) => {
      setFiles(
        acceptedFiles.map((file) =>
          Object.assign(file, {
            preview: URL.createObjectURL(file),
          })
        )
      );
    },
  });

  const openDialog = () => {
    // Note that the ref is set async,
    // so it might be null at some point
    if (dropzoneRef.current) {
      dropzoneRef.current.open();
    }
  };

  // @ts-ignore
  const thumbs = files.map(({ name, preview }: FileWithPath) => (
    <div style={thumb} key={name}>
      <div style={thumbInner}>
        <img src={preview} style={img} />
      </div>
    </div>
  ));

  const upload = () => {};

  return (
    <PageWrapper>
      <section className="container">
        <Container {...getRootProps()}>
          <input {...getInputProps()} />
          <p>
            Drag &apos;`n&apos;` drop some files here, or click to select files
          </p>
          <button type="button" onClick={openDialog}>
            Open File Dialog
          </button>
          <hr />
        </Container>
        <aside style={thumbsContainer}>{thumbs}</aside>

        <div>
          <button onClick={}></button>
        </div>
      </section>
    </PageWrapper>
  );
};

export default Upload;
