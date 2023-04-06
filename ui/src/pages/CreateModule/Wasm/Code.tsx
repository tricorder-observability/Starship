import { UploadOutlined } from '@ant-design/icons';
import { Button, Form, Upload } from 'antd';

type IProps = {
  readFileContent: (info: any) => void;
};

const Index = (props: IProps) => {
  const { readFileContent } = props;
  return (
    <Form.Item
      name="wasm"
      label={'wasm'}
      extra="WASM byte code(.wasm .wat)"
      rules={[{ required: true, message: 'Please input wasm code!' }]}
    >
      <Upload
        action=""
        maxCount={1}
        accept=".wasm,.wat"
        showUploadList={false}
        onChange={readFileContent}
      >
        <Button icon={<UploadOutlined />}>Click to upload</Button>
      </Upload>
    </Form.Item>
  );
};
export default Index;
