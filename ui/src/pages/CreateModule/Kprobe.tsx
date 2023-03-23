import { MinusCircleOutlined, PlusOutlined } from '@ant-design/icons';
import { useIntl } from '@umijs/max';
import { Button, Form, Input, Space } from 'antd';

type IProps = {};

const Index = (props: IProps) => {
  const intl = useIntl();
  return (
    <Form.Item
      wrapperCol={{ offset: 0, span: 15 }}
      label={intl.formatMessage({
        id: 'code.kprobe',
      })}
      required={true}
    >
      <Form.List
        name="probes"
        initialValue={[
          {
            target: null,
            entry: null,
            return: null,
          },
        ]}
      >
        {(fields, { add, remove }) => (
          <>
            {fields.map(({ key, name, ...restField }) => (
              <Space key={key} style={{ display: 'flex', marginBottom: 8 }} align="baseline">
                <Form.Item
                  {...restField}
                  name={[name, 'target']}
                  rules={[{ required: true, message: 'target' }]}
                >
                  <Input placeholder="target" />
                </Form.Item>
                <Form.Item
                  {...restField}
                  name={[name, 'entry']}
                  rules={[{ required: true, message: 'entry' }]}
                >
                  <Input placeholder="entry" />
                </Form.Item>
                <Form.Item
                  {...restField}
                  name={[name, 'return']}
                  rules={[{ required: true, message: 'return' }]}
                >
                  <Input placeholder="return" />
                </Form.Item>
                {fields.length > 1 && <MinusCircleOutlined onClick={() => remove(name)} />}
              </Space>
            ))}
            <Form.Item>
              <Button
                type="dashed"
                onClick={() => add()}
                block
                icon={<PlusOutlined />}
                style={{ width: '100%' }}
              >
                {intl.formatMessage({
                  id: 'code.addKprobe',
                })}
              </Button>
            </Form.Item>
          </>
        )}
      </Form.List>
    </Form.Item>
  );
};
export default Index;
