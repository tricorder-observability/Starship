import { Form, Input } from 'antd';

type IProps = {};

const Index = (props: IProps) => {
  return (
    <Form.Item
      label={'name'}
      name="name"
      rules={[
        { required: true, message: 'Please input name!' },
        { type: 'string', max: 50, message: 'Up to 50 characters!' },
      ]}
    >
      <Input style={{ width: '100%' }} />
    </Form.Item>
  );
};
export default Index;
