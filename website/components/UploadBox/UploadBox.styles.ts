import styled from "@emotion/styled";
import theme from "../../styles/colors";

export const UploadWrapper = styled.div`
  width: 40%;
  margin: 0 auto;
  padding: 10px;
  background-color: #fafafa;
  border-width: 2px;
  border-radius: 1px;
  border-color: ${theme.border};
  border-style: solid;
`;

export const Img = styled.img`
  max-width: 100%;
  height: auto;
  display: block;
  margin: auto;
`;

export const Container = styled.div`
  align-items: center;
  border-width: 2px;
  border-radius: 2px;
  border-color: #eeeeee;
  border-style: dashed;
  padding: 10px;
  background-color: #fafafa;
  color: ${theme.secondary};
  outline: none;
  text-align: center;
`;

export const FormColumn = styled.div`
  display: flex;
  flex-direction: column;
  margin-right: 10%;
`;

export const FormRow = styled.div`
  display: flex;
  flex-direction: row;
  margin-bottom: 5px;
`;

export const FormArea = styled.div`
  margin: 10px 0px;
  display: flex;
  flex-direction: row;
`;

export const FormBox = styled.div`
  display: flex;
  flex-direction: column;
  margin-bottom: 5px;
`;

export const FormLabel = styled.label<{ emphasized?: boolean }>`
  font-weight: ${({ emphasized }) => (emphasized ? `bold` : `normal`)};
  color: ${theme.secondary};
  text-transform: capitalize;
  margin-bottom: 5px;
`;

export const FormCheckBoxInput = styled.input`
  margin-right: 5px;
`;

export const FormButtonWrapper = styled.div`
  margin-left: auto;
  display: flex;
`;

export const FormButton = styled.button`
  align-self: flex-end;
  border: 0;
  border-radius: 0.25rem;
  font-size: 1rem;
  line-height: 1.2;
  white-space: nowrap;
  text-decoration: none;
  padding: 0.25rem 0.5rem;
  margin: 0.25rem;
  cursor: pointer;
  color: ${theme.primary};

  &:hover {
    background: ${theme.border};
    color: #fff;
  }
`;
