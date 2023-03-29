import { MinusCircleOutlined, PlusOutlined } from '@ant-design/icons';
import { Button, Form, Input, Select, Space } from 'antd';

type IProps = {};

const Index = (props: IProps) => {
  return (
    <Form.Item wrapperCol={{ offset: 0, span: 15 }} label={'probe'} required={true}>
      <Form.List
        name="probes"
        initialValue={[
          {
            target: null,
            entry: null,
            return: null,
            binary_path: null,
            sample_period_nanos: null,
            type: null,
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
                <Form.Item
                  {...restField}
                  name={[name, 'binary_path']}
                  rules={[{ required: true, message: 'binary_path' }]}
                >
                  <Input placeholder="binary_path" />
                </Form.Item>
                <Form.Item
                  {...restField}
                  name={[name, 'sample_period_nanos']}
                  rules={[{ required: true, message: 'sample_period_nanos' }]}
                >
                  <Input placeholder="sample_period_nanos" />
                </Form.Item>
                <Form.Item
                  {...restField}
                  name={[name, 'type']}
                  rules={[{ required: true, message: 'type' }]}
                >
                  <Select placeholder={'type'} style={{ width: 120 }}>
                    <Select.Option value={0}>0</Select.Option>
                    <Select.Option value={1}>1</Select.Option>
                  </Select>
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
                {'add probe'}
              </Button>
            </Form.Item>
          </>
        )}
      </Form.List>
    </Form.Item>
  );
};
export default Index;
