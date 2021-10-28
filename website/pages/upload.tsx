import { PageWrapper } from "components";
import type { NextPage } from "next";
import React from "react";
import { UploadBox } from "../components/UploadBox/UploadBox";

const Upload: NextPage = () => {
  return (
    <PageWrapper>
      <UploadBox />
    </PageWrapper>
  );
};

export default Upload;
