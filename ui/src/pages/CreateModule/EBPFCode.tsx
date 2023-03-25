import { useIntl } from '@umijs/max';
import { Form, Input } from 'antd';

type IProps = {};

const Index = (props: IProps) => {
  const intl = useIntl();
  return (
    <Form.Item
      label={intl.formatMessage({
        id: 'code.code',
      })}
      name="code"
      rules={[{ required: true, message: 'Please input ebpf!' }]}
    >
      <Input.TextArea rows={8} style={{ width: '100%' }} />
    </Form.Item>
  );
};
export default Index;
