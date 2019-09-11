
conda config --add channels https://mirrors.tuna.tsinghua.edu.cn/anaconda/pkgs/free/
conda config --set show_channel_urls yes
---------------------
https://www.anaconda.com/distribution/#download-section

创建一个名称为python34的虚拟环境并指定python版本为3.4(这里conda会自动找3.4中最新的版本下载)

conda  create -n python34  python=3.4

或者conda  create  --name  python34   python=3.4
————————————————
版权声明：本文为CSDN博主「代码帮」的原创文章，遵循 CC 4.0 BY-SA 版权协议，转载请附上原文出处链接及本声明。
原文链接：https://blog.csdn.net/ITLearnHall/article/details/81708148

---------------------
(https://pytorch-cn.readthedocs.io/zh/latest/notes/autograd/)
https://pytorch.org/get-started/locally/
conda install pytorch-cpu torchvision-cpu -c pytorch

conda install pytorch torchvision cudatoolkit=10.0 -c pytorch

CUDA 9.0
To install PyTorch via Anaconda, and you are using CUDA 9.0, use the following conda command:
conda install pytorch torchvision cudatoolkit=9.0 -c pytorch
CUDA 8.x
conda install pytorch torchvision cudatoolkit=8.0 -c pytorch
CUDA 10.0
conda install pytorch torchvision cudatoolkit=10.0 -c pytorch