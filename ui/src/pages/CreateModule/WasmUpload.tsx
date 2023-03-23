import { UploadOutlined } from '@ant-design/icons';
import { useIntl } from '@umijs/max';
import { Button, Form, Upload } from 'antd';

type IProps = {
  readFileContent: (info: any) => void;
};

const Index = (props: IProps) => {
  const { readFileContent } = props;
  const intl = useIntl();
  return (
    <Form.Item
      name="wasm"
      label={intl.formatMessage({
        id: 'code.wasm',
      })}
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
