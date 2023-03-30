import { UploadOutlined } from '@ant-design/icons';
import { Button, Upload } from 'antd';

type IProps = {};

const Index = (props: IProps) => {
  const readFileContent = (info: any) => {};
  return (
    <Upload action="" showUploadList={false} beforeUpload={readFileContent} directory>
      <Button icon={<UploadOutlined />}>Click to Upload</Button>
    </Upload>
  );
};
export default Index;
