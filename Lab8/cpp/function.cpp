#include <iostream>
#include <cmath>

using namespace std;

int main()
{
    float x;
     cin >> x ;

    double Sum = 1, q = 1,eps = 0.0001f;
    int n = 1;
    
    while (true)
    {
        q = q*x / n;

        Sum += q;
        n++;

        if (q < eps)
            break;

    }
    cout <<"Нахождения суммы бесконечного ряда (разложение экспоненты) с заданной точностью ε=0,0001"<< endl
         << "X=" << x << endl
         << "Summ = " << Sum << endl
         << "EXP(X)=" << exp(x);
    //}
}