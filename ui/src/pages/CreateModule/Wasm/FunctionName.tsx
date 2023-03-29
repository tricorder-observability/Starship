import { Form, Input } from 'antd';

type IProps = {};

const Index = (props: IProps) => {
  return (
    <Form.Item
      label={'fn name'}
      name="fn_name"
      rules={[{ required: true, message: 'Please input function name!' }]}
    >
      <Input style={{ width: '100%' }} />
    </Form.Item>
  );
};
export default Index;
